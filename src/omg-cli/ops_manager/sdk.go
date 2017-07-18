package ops_manager

import (
	"crypto/tls"
	"log"
	"os"
	"time"

	"encoding/json"

	"fmt"

	"net/http"

	"bytes"

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
	username              string
	password              string
	logger                *log.Logger
	unauthenticatedClient network.UnauthenticatedClient
	client                network.OAuthClient
	httpClient            *http.Client
}

func NewSdk(target, username, password string, skipSSLValidation bool, logger log.Logger) (*Sdk, error) {
	client, err := network.NewOAuthClient(target, username, password, "", "", skipSSLValidation, true, time.Duration(requestTimeout)*time.Second)
	if err != nil {
		return nil, err
	}

	logger.SetPrefix(fmt.Sprintf("%s[OM SDK] ", logger.Prefix()))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipSSLValidation},
	}

	return &Sdk{target: target,
		username:              username,
		password:              password,
		logger:                &logger,
		unauthenticatedClient: network.NewUnauthenticatedClient(target, skipSSLValidation, time.Duration(requestTimeout)*time.Second),
		client:                client,
		httpClient:            &http.Client{Transport: tr},
	}, nil
}

func (om *Sdk) SetupAuth(decryptionPhrase string) error {
	setupService := api.NewSetupService(om.unauthenticatedClient)

	cmd := commands.NewConfigureAuthentication(setupService, om.logger)
	return cmd.Execute([]string{
		"--username", om.username,
		"--password", om.password,
		"--decryption-passphrase", decryptionPhrase})
}

type UnlockRequest struct {
	Passphrase string `json:"passphrase"`
}

func (om *Sdk) Unlock(decryptionPhrase string) error {
	om.logger.Println("decrypting Ops Manager")
	unlockReq := UnlockRequest{decryptionPhrase}
	body, err := json.Marshal(&unlockReq)

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v0/unlock", om.target), bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	_, err = om.httpClient.Do(req)

	return err
}

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

func (om *Sdk) ApplyChanges() error {
	installationsService := api.NewInstallationsService(om.client)
	logWriter := commands.NewLogWriter(os.Stdout)
	cmd := commands.NewApplyChanges(installationsService, logWriter, om.logger, poolingIntervalSec)

	return cmd.Execute(nil)
}

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

func (om *Sdk) StageProduct(name, version string) error {
	diagnosticService := api.NewDiagnosticService(om.client)
	availableProductsService := api.NewAvailableProductsService(om.client, progress.NewBar(), uilive.New())
	stagedProductsService := api.NewStagedProductsService(om.client)
	cmd := commands.NewStageProduct(stagedProductsService, availableProductsService, diagnosticService, om.logger)
	return cmd.Execute([]string{
		"--product-name", name,
		"--product-version", version,
	})
}

func (om *Sdk) Ready() bool {
	om.logger.Print("checking if Ops Manager is ready... ")

	req, err := http.NewRequest("GET", om.target, nil)
	if err != nil {
		return false
	}
	resp, err := om.httpClient.Do(req)

	om.logger.Printf("got: %d\n", resp.StatusCode)

	return resp.StatusCode < 500
}

func (om *Sdk) AvaliableProducts() ([]api.ProductInfo, error) {
	service := api.NewAvailableProductsService(om.client, progress.NewBar(), uilive.New())
	out, err := service.List()
	if err != nil {
		return nil, err
	}

	return out.ProductsList, nil
}

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
