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
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"io"

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
	requestTimeout     = 1800
	poolingIntervalSec = 10
)

type Sdk struct {
	target     string
	creds      config.OpsManagerCredentials
	logger     *log.Logger
	api        api.Api
	httpClient *http.Client
}

// NewSdk creates an authenticated session and object to interact with Ops Manager
func NewSdk(target string, creds config.OpsManagerCredentials, logger log.Logger) (*Sdk, error) {
	timeout := time.Duration(requestTimeout) * time.Second
	authedClient, err := network.NewOAuthClient(target, creds.Username, creds.Password, "", "", creds.SkipSSLVerification, true, timeout, timeout)
	if err != nil {
		return nil, err
	}

	unauthenticatedClient := network.NewUnauthenticatedClient(target, creds.SkipSSLVerification, timeout, timeout)

	liveWriter := uilive.New()
	liveWriter.Out = os.Stderr
	authedProgressClient := network.NewProgressClient(authedClient, progress.NewBar(), liveWriter)
	unauthenticatedProgressClient := network.NewProgressClient(unauthenticatedClient, progress.NewBar(), liveWriter)

	api := api.New(api.ApiInput{
		Client:                 authedClient,
		UnauthedClient:         unauthenticatedClient,
		ProgressClient:         authedProgressClient,
		UnauthedProgressClient: unauthenticatedProgressClient,
		Logger:                 &logger,
	})

	logger.SetPrefix(fmt.Sprintf("%s[OM SDK] ", logger.Prefix()))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: creds.SkipSSLVerification},
	}

	return &Sdk{target: target,
		creds:  creds,
		logger: &logger,

		api:        api,
		httpClient: &http.Client{Transport: tr},
	}, nil
}

// SetupAuth configures the initial username, password, and decryptionPhrase
func (om *Sdk) SetupAuth() error {
	cmd := commands.NewConfigureAuthentication(om.api, om.logger)
	return cmd.Execute([]string{
		"--username", om.creds.Username,
		"--password", om.creds.Password,
		"--decryption-passphrase", om.creds.DecryptionPhrase})
}

// Unlock decrypts Ops Manager. This is needed after a reboot before attempting to authenticate.
// This task runs asynchronously. Query the status by invoking ReadyForAuth.
func (om *Sdk) Unlock() error {
	om.logger.Println("decrypting Ops Manager")
	unlockReq := UnlockRequest{om.creds.DecryptionPhrase}
	body, err := json.Marshal(&unlockReq)

	req, err := om.newRequest("PUT", fmt.Sprintf("%s/api/v0/unlock", om.target), bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := om.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// ReadyForAuth checks if the Ops Manager authentication system is ready
func (om *Sdk) ReadyForAuth() bool {
	req, err := om.newRequest("GET", fmt.Sprintf("%s/login/ensure_availability", om.target), nil)
	if err != nil {
		return false
	}
	resp, err := om.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// When OpsMan is online/decrypted it redirects its auth system. UAA is expected for OMG.
	return resp.StatusCode == 200 && strings.Contains(resp.Request.URL.Path, "/uaa/login")
}

// SetupBosh applies the provided configuration to the BOSH director tile
func (om *Sdk) SetupBosh(iaas commands.GCPIaaSConfiguration, director commands.DirectorConfiguration, azs commands.AvailabilityZonesConfiguration, networks commands.NetworksConfiguration, networkAssignment commands.NetworkAssignment, resources commands.ResourceConfiguration) error {

	cmd := commands.NewConfigureDirector(os.Environ, om.api, om.logger)

	iaasBytes, err := json.Marshal(iaas)
	if err != nil {
		return err
	}

	directorBytes, err := json.Marshal(director)
	if err != nil {
		return err
	}

	azBytes, err := json.Marshal(azs)
	if err != nil {
		return err
	}

	networksBytes, err := json.Marshal(networks)
	if err != nil {
		return err
	}

	networkAssignmentBytes, err := json.Marshal(networkAssignment)
	if err != nil {
		return err
	}

	resourceBytes, err := json.Marshal(resources)
	if err != nil {
		return err
	}

	return cmd.Execute([]string{
		"--iaas-configuration", string(iaasBytes),
		"--director-configuration", string(directorBytes),
		"--az-configuration", string(azBytes),
		"--networks-configuration", string(networksBytes),
		"--network-assignment", string(networkAssignmentBytes),
		"--resource-configuration", string(resourceBytes)})
}

// ApplyChanges deploys pending changes to Ops Manager
func (om *Sdk) ApplyChanges() error {
	logWriter := commands.NewLogWriter(os.Stdout)
	cmd := commands.NewApplyChanges(om.api, om.api, logWriter, om.logger, poolingIntervalSec)

	return cmd.Execute(nil)
}

func (om *Sdk) ApplyDirector() error {
	logWriter := commands.NewLogWriter(os.Stdout)
	cmd := commands.NewApplyChanges(om.api, om.api, logWriter, om.logger, poolingIntervalSec)
	return cmd.Execute([]string{"--skip-deploy-products"})
}

// UploadProduct pushes a given file located locally at path to the target
func (om *Sdk) UploadProduct(path string) error {
	form := formcontent.NewForm()
	cmd := commands.NewUploadProduct(form, extractor.MetadataExtractor{}, om.api, om.logger)

	return cmd.Execute([]string{
		"--product", path})
}

// UploadStemcell pushes a given stemcell located locally at path to the target
func (om *Sdk) UploadStemcell(path string) error {
	form := formcontent.NewForm()
	cmd := commands.NewUploadStemcell(form, om.api, om.logger)

	return cmd.Execute([]string{
		"--stemcell", path})
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
	req, err := om.newRequest("GET", om.target, nil)
	if err != nil {
		return false
	}
	resp, err := om.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode < 500
}

// AvaliableProducts lists products that are uploaded to Ops Manager.
func (om *Sdk) AvaliableProducts() ([]api.ProductInfo, error) {
	out, err := om.api.ListAvailableProducts()
	if err != nil {
		return nil, err
	}

	return out.ProductsList, nil
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

	resp, err := om.curl(fmt.Sprintf("api/v0/staged/products/%s/properties", productGuid), "GET", nil)
	if err != nil {
		return nil, err
	}

	var prop ProductProperties
	if err := json.Unmarshal(resp, &prop); err != nil {
		return nil, err
	}

	return &prop, nil
}

// GetDirector fetches settings for the BOSH director
func (om *Sdk) GetDirector() (*DirectorProperties, error) {
	resp, err := om.curl("/api/v0/staged/director/properties", "GET", nil)
	if err != nil {
		return nil, err
	}

	var prop DirectorProperties
	if err := json.Unmarshal(resp, &prop); err != nil {
		return nil, err
	}

	return &prop, nil
}

// GetResource fetches resource settings for a specific job of a tile
func (om *Sdk) GetResource(tileName, jobName string) (*Resource, error) {
	productGuid, err := om.productGuidByType(tileName)
	if err != nil {
		return nil, err
	}

	jobGuid, err := om.jobGuidByName(productGuid, jobName)
	if err != nil {
		return nil, err
	}

	resp, err := om.curl(fmt.Sprintf("/api/v0/staged/products/%s/jobs/%s/resource_config", productGuid, jobGuid), "GET", nil)
	if err != nil {
		return nil, err
	}

	var prop Resource
	if err := json.Unmarshal(resp, &prop); err != nil {
		return nil, err
	}

	return &prop, nil
}

func (om *Sdk) curl(path, method string, data io.Reader) ([]byte, error) {
	resp, err := om.api.Curl(api.RequestServiceCurlInput{
		Path:   path,
		Method: method,
		Data:   data,
	})

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check if the OpsMan API returned an error
	errResp := ErrorResponse{make(map[string][]string)}
	if err := json.Unmarshal(body, &errResp); err == nil {
		if len(errResp.Errors) != 0 {
			return nil, fmt.Errorf("error from Ops Manager API requesting %s: %v", path, errResp.Errors)
		}
	}

	return body, nil
}

func (om *Sdk) getProducts() ([]Product, error) {
	body, err := om.curl("api/v0/deployed/products", http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	var resp []Product
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("malformed products response: %s", string(body))
	}

	return resp, nil
}

func (om *Sdk) productGuidByType(product string) (string, error) {
	products, err := om.getProducts()
	if err != nil {
		return "", err
	}

	appGuid := ""
	for _, p := range products {
		if p.Type == product {
			appGuid = p.Guid
			break
		}
	}

	if appGuid == "" {
		return "", fmt.Errorf("could not find installed application by name: %s", product)
	}

	return appGuid, nil
}

func (om *Sdk) jobGuidByName(productGuid, jobName string) (string, error) {
	resp, err := om.curl(fmt.Sprintf("/api/v0/staged/products/%s/jobs", productGuid), "GET", nil)
	if err != nil {
		return "", err
	}

	var jobResp JobsResponse
	if err := json.Unmarshal(resp, &jobResp); err != nil {
		return "", err
	}

	jobGuid := ""
	for _, j := range jobResp.Jobs {
		if j.Name == jobName {
			jobGuid = j.Guid
			break
		}
	}

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
	return om.getCredential(fmt.Sprintf("api/v0/deployed/products/%s/credentials/%s", productGuid, credential))
}

func (om *Sdk) GetDirectorCredentials(credential string) (*SimpleCredential, error) {
	return om.getCredential(fmt.Sprintf("api/v0/deployed/director/credentials/%s", credential))
}

func (om *Sdk) getCredential(path string) (*SimpleCredential, error) {
	body, err := om.curl(path, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	var resp CredentialResponse
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
	body, err := om.curl(fmt.Sprintf("api/v0/deployed/products/%s/static_ips", boshGuid), http.MethodGet, nil)
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
	cmd := commands.NewDeleteInstallation(om.api, logWriter, om.logger, poolingIntervalSec)

	return cmd.Execute(nil)
}

func (om *Sdk) newRequest(method, url string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, body)
	if req != nil {
		req.Header.Set("User-Agent", version.UserAgent())
	}
	return
}
