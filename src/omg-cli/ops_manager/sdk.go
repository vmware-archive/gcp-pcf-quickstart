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
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"omg-cli/config"

	"io"

	"io/ioutil"

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
	target                string
	creds                 config.OpsManagerCredentials
	logger                *log.Logger
	unauthenticatedClient network.UnauthenticatedClient
	client                network.OAuthClient
	httpClient            *http.Client
}

type CredentialResponse struct {
	Credential CredentialWrapper `json:"credential"`
}

type CredentialWrapper struct {
	Type  string           `json:"type"`
	Value SimpleCredential `json:"value"`
}

type SimpleCredential struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

func NewSdk(target string, creds config.OpsManagerCredentials, logger log.Logger) (*Sdk, error) {
	client, err := network.NewOAuthClient(target, creds.Username, creds.Password, "", "", creds.SkipSSLVerification, true, time.Duration(requestTimeout)*time.Second)
	if err != nil {
		return nil, err
	}

	logger.SetPrefix(fmt.Sprintf("%s[OM SDK] ", logger.Prefix()))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: creds.SkipSSLVerification},
	}

	return &Sdk{target: target,
		creds:                 creds,
		logger:                &logger,
		unauthenticatedClient: network.NewUnauthenticatedClient(target, creds.SkipSSLVerification, time.Duration(requestTimeout)*time.Second),
		client:                client,
		httpClient:            &http.Client{Transport: tr},
	}, nil
}

// SetupAuth configures the initial username, password, and decryptionPhrase
func (om *Sdk) SetupAuth() error {
	setupService := api.NewSetupService(om.unauthenticatedClient)

	cmd := commands.NewConfigureAuthentication(setupService, om.logger)
	return cmd.Execute([]string{
		"--username", om.creds.Username,
		"--password", om.creds.Password,
		"--decryption-passphrase", om.creds.DecryptionPhrase})
}

type UnlockRequest struct {
	Passphrase string `json:"passphrase"`
}

// Unlock decrypts Ops Manager. This is needed after a reboot before attempting to authenticate.
// This task runs asynchronusly. Query the status by invoking ReadyForAuth.
func (om *Sdk) Unlock() error {
	om.logger.Println("decrypting Ops Manager")
	unlockReq := UnlockRequest{om.creds.DecryptionPhrase}
	body, err := json.Marshal(&unlockReq)

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v0/unlock", om.target), bytes.NewReader(body))
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/login/ensure_availability", om.target), nil)
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
	boshService := api.NewBoshFormService(om.client)
	diagnosticService := api.NewDiagnosticService(om.client)
	cmd := commands.NewConfigureBosh(boshService, diagnosticService, om.logger)

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
	installationsService := api.NewInstallationsService(om.client)
	logWriter := commands.NewLogWriter(os.Stdout)
	cmd := commands.NewApplyChanges(installationsService, logWriter, om.logger, poolingIntervalSec)

	return cmd.Execute(nil)
}

// UploadProduct pushes a given file located locally at path to the target
func (om *Sdk) UploadProduct(path string) error {
	liveWriter := uilive.New()
	availableProductsService := api.NewAvailableProductsService(om.client, progress.NewBar(), liveWriter)

	form, err := formcontent.NewForm()
	if err != nil {
		return err
	}

	cmd := commands.NewUploadProduct(form, extractor.ProductUnzipper{}, availableProductsService, om.logger)

	return cmd.Execute([]string{
		"--product", path})
}

// UploadStemcell pushes a given stemcell located locally at path to the target
func (om *Sdk) UploadStemcell(path string) error {
	diagnosticService := api.NewDiagnosticService(om.client)
	form, err := formcontent.NewForm()
	if err != nil {
		return err
	}

	uploadStemcellService := api.NewUploadStemcellService(om.client, progress.NewBar())
	cmd := commands.NewUploadStemcell(form, uploadStemcellService, diagnosticService, om.logger)

	return cmd.Execute([]string{
		"--stemcell", path})
}

// StageProduct moves a given name, version to the list of tiles that will be deployed
func (om *Sdk) StageProduct(tile config.OpsManagerMetadata) error {
	diagnosticService := api.NewDiagnosticService(om.client)
	availableProductsService := api.NewAvailableProductsService(om.client, progress.NewBar(), uilive.New())
	stagedProductsService := api.NewStagedProductsService(om.client)
	cmd := commands.NewStageProduct(stagedProductsService, availableProductsService, diagnosticService, om.logger)
	return cmd.Execute([]string{
		"--product-name", tile.Name,
		"--product-version", tile.Version,
	})
}

// Online checks if Ops Manager is running on the target.
func (om *Sdk) Online() bool {
	req, err := http.NewRequest("GET", om.target, nil)
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
	service := api.NewAvailableProductsService(om.client, progress.NewBar(), uilive.New())
	out, err := service.List()
	if err != nil {
		return nil, err
	}

	return out.ProductsList, nil
}

// ConfigureProduct sets up the settings for a given tile by name
func (om *Sdk) ConfigureProduct(name, networks, properties string, resources string) error {
	stagedProductsService := api.NewStagedProductsService(om.client)
	jobsService := api.NewJobsService(om.client)
	cmd := commands.NewConfigureProduct(stagedProductsService, jobsService, om.logger)

	return cmd.Execute([]string{
		"--product-name", name,
		"--product-network", networks,
		"--product-properties", properties,
		"--product-resources", resources,
	})
}

type ErrorResponse struct {
	Errors map[string][]string `json:errors`
}

func (om *Sdk) curl(path, method string, data io.Reader) ([]byte, error) {
	request, err := http.NewRequest(method, fmt.Sprintf("%s/%s", om.target, path), data)
	if err != nil {
		return nil, err
	}
	resp, err := om.client.Do(request)

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

func (om *Sdk) GetCredentials(productGuid, credential string) (*SimpleCredential, error) {
	body, err := om.curl(fmt.Sprintf("api/v0/deployed/products/%s/credentials/%s", productGuid, credential), http.MethodGet, nil)
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
