package ops_manager

import (
	"log"
	"os"
	"time"

	"encoding/json"

	"github.com/pivotal-cf/om/api"
	"github.com/pivotal-cf/om/commands"
	"github.com/pivotal-cf/om/network"
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

func NewSdk(target, username, password string, skipSSLValidation bool) (*Sdk, error) {
	client, err := network.NewOAuthClient(target, username, password, "", "", skipSSLValidation, true, time.Duration(requestTimeout)*time.Second)
	if err != nil {
		return nil, err
	}

	return &Sdk{target: target,
		username:              username,
		password:              password,
		skipSSLValidation:     skipSSLValidation,
		logger:                log.New(os.Stdout, "[OM SDK] ", 0),
		unauthenticatedClient: network.NewUnauthenticatedClient(target, skipSSLValidation, time.Duration(requestTimeout)*time.Second),
		client:                client,
	}, nil
}

func (om *Sdk) SetupAuth(decryptionPhrase string) error {
	setupService := api.NewSetupService(om.unauthenticatedClient)

	cmd := commands.NewConfigureAuthentication(setupService, om.logger)
	return cmd.Execute([]string{"--username", om.username, "--password", om.password, "--decryption-passphrase", decryptionPhrase})
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
