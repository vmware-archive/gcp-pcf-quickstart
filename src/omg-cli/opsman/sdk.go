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

package opsman

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"omg-cli/config"

	"github.com/gosuri/uilive"
	"github.com/pivotal-cf/om/api"
	"github.com/pivotal-cf/om/network"
	"github.com/pivotal-cf/om/progress"
)

const (
	connectTimeout     = 5
	requestTimeout     = 1800
	pollingIntervalSec = 10
)

// Sdk interacts with the Ops Manager's API.
type Sdk struct {
	unauthenticatedClient network.UnauthenticatedClient
	client                network.OAuthClient
	api                   api.Api
	creds                 config.OpsManagerCredentials
	logger                *log.Logger
	target                string
}

// NewSdk creates an authenticated session and object to interact with Ops Manager
func NewSdk(target string, creds config.OpsManagerCredentials, logger *log.Logger) (*Sdk, error) {
	client, err := network.NewOAuthClient(target, creds.Username, creds.Password, "", "",
		creds.SkipSSLVerification, true, time.Duration(requestTimeout)*time.Second, time.Duration(connectTimeout)*time.Second)
	if err != nil {
		return nil, err
	}

	unauthenticatedClient := network.NewUnauthenticatedClient(target, creds.SkipSSLVerification,
		time.Duration(requestTimeout)*time.Second,
		time.Duration(connectTimeout)*time.Second)

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
			Logger:                 logger,
		}),
		creds:  creds,
		logger: logger,
		target: target,
	}
	return sdk, nil
}

func (om *Sdk) getProducts() ([]api.DeployedProductOutput, error) {
	return om.api.ListDeployedProducts()
}

func (om *Sdk) productGUIDByType(product string) (string, error) {
	products, err := om.getProducts()
	if err != nil {
		return "", err
	}

	appGUID := ""
	for _, p := range products {
		if p.Type == product {
			appGUID = p.GUID
			break
		}
	}

	if appGUID == "" {
		return "", fmt.Errorf("could not find installed application by name: %s", product)
	}

	return appGUID, nil
}

func (om *Sdk) jobGUIDByName(productGUID, jobName string) (string, error) {
	jobs, err := om.api.ListStagedProductJobs(productGUID)
	if err != nil {
		return "", err
	}

	jobGUID := jobs[jobName]
	if jobGUID == "" {
		return "", fmt.Errorf("Job %s not found for product %s", jobName, productGUID)
	}

	return jobGUID, nil
}

// GetCredentials returns a credential by name.
func (om *Sdk) GetCredentials(name, credential string) (*SimpleCredential, error) {
	productGUID, err := om.productGUIDByType(name)
	if err != nil {
		return nil, err
	}
	out, err := om.api.GetDeployedProductCredential(api.GetDeployedProductCredentialInput{
		DeployedGUID:        productGUID,
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

// GetDirectorCredentials returns the BOSH Director's credentials.
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
		return nil, fmt.Errorf("received an empty credential: %s", string(body))
	}

	return &resp.Credential.Value, nil
}

// GetDirectorIP returns the IP address of the BOSH Director.
func (om *Sdk) GetDirectorIP() (string, error) {
	boshGUID, err := om.productGUIDByType("p-bosh")
	if err != nil {
		return "", err
	}
	out, err := om.api.Curl(api.RequestServiceCurlInput{
		Path:   fmt.Sprintf("api/v0/deployed/products/%s/static_ips", boshGUID),
		Method: http.MethodGet,
	})
	if err != nil {
		return "", err
	}

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
	return "", errors.New("static_ips response had no director Job")
}
