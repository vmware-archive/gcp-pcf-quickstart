package ops_manager

import (
	"log"
	"os"
	"time"

	"encoding/json"

	"fmt"

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
	skipSSLValidation     bool
	logger                *log.Logger
	unauthenticatedClient network.UnauthenticatedClient
	client                network.OAuthClient
}

func NewSdk(target, username, password string, skipSSLValidation bool, logger log.Logger) (*Sdk, error) {
	client, err := network.NewOAuthClient(target, username, password, "", "", skipSSLValidation, true, time.Duration(requestTimeout)*time.Second)
	if err != nil {
		return nil, err
	}

	logger.SetPrefix(fmt.Sprintf("%s[OM SDK] ", logger.Prefix()))

	return &Sdk{target:        target,
		username:              username,
		password:              password,
		skipSSLValidation:     skipSSLValidation,
		logger:                &logger,
		unauthenticatedClient: network.NewUnauthenticatedClient(target, skipSSLValidation, time.Duration(requestTimeout)*time.Second),
		client:                client,
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

func (om *Sdk) SetupBosh(iaas commands.GCPIaaSConfiguration, director commands.DirectorConfiguration, azs commands.AvailabilityZonesConfiguration, networks commands.NetworksConfiguration, networkAssignment commands.NetworkAssignment) error {
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

	return cmd.Execute([]string{
		"--iaas-configuration", string(iaasBytes),
		"--director-configuration", string(directorBytes),
		"--az-configuration", string(azBytes),
		"--networks-configuration", string(networksBytes),
		"--network-assignment", string(networkAssignmentBytes)})
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

func (om *Sdk) AvaliableProducts() ([]api.ProductInfo, error) {
	service := api.NewAvailableProductsService(om.client, progress.NewBar(), uilive.New())
	out, err := service.List()
	if err != nil {
		return nil, err
	}

	return out.ProductsList, nil
}

func (om *Sdk) ConfigureProduct(name, networks, properties string) error {
	stagedProductsService := api.NewStagedProductsService(om.client)
	jobsService := api.NewJobsService(om.client)
	cmd := commands.NewConfigureProduct(stagedProductsService, jobsService, om.logger)

	return cmd.Execute([]string{
		"--product-name", name,
		"--product-network", networks,
		"--product-properties", properties,
	})
}
