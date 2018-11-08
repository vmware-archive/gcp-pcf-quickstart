/*
 * Copyright 2017 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ops_manager

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"omg-cli/config"
	"omg-cli/version"

	"github.com/gosuri/uilive"
	"github.com/pivotal-cf/om/api"
	"github.com/pivotal-cf/om/commands"
	"github.com/pivotal-cf/om/extractor"
	"github.com/pivotal-cf/om/formcontent"
	"github.com/pivotal-cf/om/network"
	"github.com/pivotal-cf/om/progress"
)

const (
	connectTimeout     = 5
	requestTimeout     = 1800
	pollingIntervalSec = 10
)

type Sdk struct {
	unauthenticatedClient network.UnauthenticatedClient
	client                network.OAuthClient
	api                   api.Api
	creds                 config.OpsManagerCredentials
	logger                *log.Logger
}

// NewSdk creates an authenticated session and object to interact with Ops Manager
func NewSdk(target string, creds config.OpsManagerCredentials, logger log.Logger) (*Sdk, error) {
	client, err := network.NewOAuthClient(target, creds.Username, creds.Password, "", "",
		creds.SkipSSLVerification, true, time.Duration(requestTimeout)*time.Second, time.Duration(connectTimeout)*time.Second)
	unauthenticatedClient := network.NewUnauthenticatedClient(target, creds.SkipSSLVerification,
		time.Duration(requestTimeout)*time.Second,
		time.Duration(connectTimeout)*time.Second)
	if err != nil {
		return nil, err
	}

	logger.SetPrefix(fmt.Sprintf("%s[OM SDK] ", logger.Prefix()))

	live := uilive.New()
	live.Out = os.Stderr

	sdk := &Sdk{
		client:                client,
		unauthenticatedClient: unauthenticatedClient,
		api: api.New(api.ApiInput{
			Client:                 client,
			UnauthedClient:         unauthenticatedClient,
			ProgressClient:         network.NewProgressClient(client, progress.NewBar(), live),
			UnauthedProgressClient: network.NewProgressClient(unauthenticatedClient, progress.NewBar(), live),
			Logger:                 &logger,
		}),
		creds:  creds,
		logger: &logger,
	}
	return sdk, nil
}

// SetupAuth configures the initial username, password, and decryptionPhrase
func (om *Sdk) SetupAuth() error {
	availability, err := om.api.EnsureAvailability(api.EnsureAvailabilityInput{})
	if err != nil {
		return fmt.Errorf("could not determine initial auth configuration status: %v", err)
	}

	if availability.Status == api.EnsureAvailabilityStatusUnknown {
		return errors.New("could not determine initial auth configuration status: unexpected status")
	}
	if availability.Status != api.EnsureAvailabilityStatusUnstarted {
		om.logger.Printf("configuration previously completed, skipping configuration")
		return nil
	}
	om.logger.Println("configuring internal userstore...")
	_, err = om.api.Setup(api.SetupInput{
		IdentityProvider:                 "internal",
		AdminUserName:                    om.creds.Username,
		AdminPassword:                    om.creds.Password,
		AdminPasswordConfirmation:        om.creds.Password,
		DecryptionPassphrase:             om.creds.DecryptionPhrase,
		DecryptionPassphraseConfirmation: om.creds.DecryptionPhrase,
		EULAAccepted:                     "true",
	})
	if err != nil {
		return fmt.Errorf("could not configure auth: %v", err)
	}
	for availability.Status != api.EnsureAvailabilityStatusComplete {
		availability, err = om.api.EnsureAvailability(api.EnsureAvailabilityInput{})
		if err != nil {
			return fmt.Errorf("could not determine final auth configuration status: %v", err)
		}
	}
	om.logger.Println("auth configuration complete")
	return nil
}

// Unlock decrypts Ops Manager. This is needed after a reboot before attempting to authenticate.
// This task runs asynchronously. Query the status by invoking ReadyForAuth.
func (om *Sdk) Unlock() error {
	om.logger.Println("decrypting Ops Manager")

	unlockReq := UnlockRequest{om.creds.DecryptionPhrase}
	body, err := json.Marshal(&unlockReq)

	_, err = om.api.Curl(api.RequestServiceCurlInput{
		Path:   "/api/v0/unlock",
		Method: "PUT",
		Data:   bytes.NewReader(body),
	})

	return err
}

// ReadyForAuth checks if the Ops Manager authentication system is ready
func (om *Sdk) ReadyForAuth() bool {
	resp, err := om.api.EnsureAvailability(api.EnsureAvailabilityInput{})
	return err == nil && resp.Status == api.EnsureAvailabilityStatusComplete
}

// SetupBosh applies the provided configuration to the BOSH director tile
func (om *Sdk) SetupBosh(configYML []byte) error {
	f, err := ioutil.TempFile("", "director-config")
	if err != nil {
		return fmt.Errorf("cannot create temp file for director configuration: %v", err)
	}
	defer os.Remove(f.Name())
	_, err = f.Write(configYML)
	if err != nil {
		return fmt.Errorf("cannot write to director yaml file: %v", err)
	}
	f.Close()

	cmd := commands.NewConfigureDirector(os.Environ, om.api, om.logger)
	return cmd.Execute([]string{"--config", f.Name()})
}

// ApplyChanges deploys pending changes to Ops Manager
func (om *Sdk) ApplyChanges(args []string) error {
	logWriter := commands.NewLogWriter(os.Stdout)
	cmd := commands.NewApplyChanges(om.api, om.api, logWriter, om.logger, 10)
	return cmd.Execute(args)
}

func (om *Sdk) ApplyDirector() error {
	logWriter := commands.NewLogWriter(os.Stdout)
	cmd := commands.NewApplyChanges(om.api, om.api, logWriter, om.logger, 10)
	return cmd.Execute([]string{"--skip-deploy-products"})
}

// UploadProduct pushes a given file located locally at path to the target
func (om *Sdk) UploadProduct(path string) error {
	form := formcontent.NewForm()
	cmd := commands.NewUploadProduct(form, extractor.MetadataExtractor{}, om.api, om.logger)
	return cmd.Execute([]string{"--product", path})
}

// UploadStemcell pushes a given stemcell located locally at path to the target
func (om *Sdk) UploadStemcell(path string) error {
	form := formcontent.NewForm()
	cmd := commands.NewUploadStemcell(form, om.api, om.logger)
	return cmd.Execute([]string{"--stemcell", path})
}

// StageProduct moves a given name, version to the list of tiles that will be deployed
func (om *Sdk) StageProduct(tile config.OpsManagerMetadata) error {
	cmd := commands.NewStageProduct(om.api, om.logger)
	return cmd.Execute([]string{
		"--product-name", tile.Name,
		"--product-version", tile.Version,
	})
}

// Online checks if Ops Manager is running on the target.
func (om *Sdk) Online() bool {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	resp, err := om.unauthenticatedClient.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode < 500
}

// AvaliableProducts lists products that are uploaded to Ops Manager.
func (om *Sdk) AvaliableProducts() ([]api.ProductInfo, error) {
	products, err := om.api.ListAvailableProducts()
	if err != nil {
		return nil, err
	}

	return products.ProductsList, nil
}

// ConfigureProduct sets up the settings for a given tile by name
func (om *Sdk) ConfigureProduct(name, networks, properties string, resources string) error {
	cmd := commands.NewConfigureProduct(os.Environ, om.api, om.logger)
	return cmd.Execute([]string{
		"--product-name", name,
		"--product-network", networks,
		"--product-properties", properties,
		"--product-resources", resources,
	})
}

// GetProduct fetches settings for a given tile by name
func (om *Sdk) GetProduct(name string) (*ProductProperties, error) {
	productGuid, err := om.productGuidByType(name)
	if err != nil {
		return nil, err
	}

	props, err := om.api.GetStagedProductProperties(productGuid)
	if err != nil {
		return nil, err
	}

	return &ProductProperties{
		Properties: props,
	}, nil
}

// GetDirector fetches settings for the BOSH director
func (om *Sdk) GetDirector() (map[string]map[string]interface{}, error) {
	props, err := om.api.GetStagedDirectorProperties()
	if err != nil {
		return nil, err
	}

	return props, nil
}

// GetResource fetches resource settings for a specific job of a tile
func (om *Sdk) GetResource(tileName, jobName string) (*api.JobProperties, error) {
	productGuid, err := om.productGuidByType(tileName)
	if err != nil {
		return nil, err
	}

	jobGuid, err := om.jobGuidByName(productGuid, jobName)
	if err != nil {
		return nil, err
	}

	props, err := om.api.GetStagedProductJobResourceConfig(productGuid, jobGuid)
	if err != nil {
		return nil, err
	}
	return &props, nil
}

func (om *Sdk) getProducts() ([]api.DeployedProductOutput, error) {
	return om.api.ListDeployedProducts()
}

func (om *Sdk) productGuidByType(product string) (string, error) {
	products, err := om.getProducts()
	if err != nil {
		return "", err
	}

	appGuid := ""
	for _, p := range products {
		if p.Type == product {
			appGuid = p.GUID
			break
		}
	}

	if appGuid == "" {
		return "", fmt.Errorf("could not find installed application by name: %s", product)
	}

	return appGuid, nil
}

func (om *Sdk) jobGuidByName(productGuid, jobName string) (string, error) {
	jobs, err := om.api.ListStagedProductJobs(productGuid)
	if err != nil {
		return "", err
	}

	jobGuid := jobs[jobName]
	if jobGuid == "" {
		return "", fmt.Errorf("job %s not found for product %s", jobName, productGuid)
	}

	return jobGuid, nil
}

func (om *Sdk) GetCredentials(name, credential string) (*SimpleCredential, error) {
	productGuid, err := om.productGuidByType(name)
	if err != nil {
		return nil, err
	}
	out, err := om.api.GetDeployedProductCredential(api.GetDeployedProductCredentialInput{
		DeployedGUID:        productGuid,
		CredentialReference: credential,
	})
	if err != nil {
		return nil, err
	}
	return &SimpleCredential{
		Identity: out.Credential.Value["identity"],
		Password: out.Credential.Value["password"],
	}, nil
}

func (om *Sdk) GetDirectorCredentials(credential string) (*SimpleCredential, error) {
	return om.getCredential(fmt.Sprintf("api/v0/deployed/director/credentials/%s", credential))
}

func (om *Sdk) getCredential(path string) (*SimpleCredential, error) {
	out, err := om.api.Curl(api.RequestServiceCurlInput{
		Path:   path,
		Method: http.MethodGet,
	})
	if err != nil {
		return nil, err
	}

	var resp CredentialResponse
	body, err := ioutil.ReadAll(out.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("malformed credentials response: %s", string(body))
	}

	if resp.Credential.Value.Password == "" || resp.Credential.Value.Identity == "" {
		return nil, fmt.Errorf("recieved an empty credential: %s", string(body))
	}

	return &resp.Credential.Value, nil
}

func (om *Sdk) GetDirectorIP() (string, error) {
	boshGuid, err := om.productGuidByType("p-bosh")
	if err != nil {
		return "", err
	}
	out, err := om.api.Curl(api.RequestServiceCurlInput{
		Path:   fmt.Sprintf("api/v0/deployed/products/%s/static_ips", boshGuid),
		Method: http.MethodGet,
	})
	body, err := ioutil.ReadAll(out.Body)
	if err != nil {
		return "", err
	}

	var boshIPs []StaticIP
	if err := json.Unmarshal(body, &boshIPs); err != nil {
		return "", fmt.Errorf("malformed static_ips response: %s", string(body))
	}
	for _, ip := range boshIPs {
		if strings.HasPrefix(ip.Name, "director") {
			return ip.IPs[0], nil
		}
	}
	return "", errors.New("static_ips response had no director job")
}

func (om *Sdk) DeleteInstallation() error {
	logWriter := commands.NewLogWriter(os.Stdout)
	cmd := commands.NewDeleteInstallation(om.api, logWriter, om.logger, pollingIntervalSec)
	return cmd.Execute(nil)
}

func (om *Sdk) newRequest(method, url string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, body)
	if req != nil {
		req.Header.Set("User-Agent", version.UserAgent())
	}
	return
}
