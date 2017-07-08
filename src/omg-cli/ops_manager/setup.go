package ops_manager

import (
	"fmt"
	"log"
	"omg-cli/config"
	"os"
	"time"

	"github.com/pivotal-cf/om/api"
	"github.com/pivotal-cf/om/commands"
	"github.com/pivotal-cf/om/network"
)

const requestTimeout = 1800

type OpsManager struct {
	cfg *config.Config
}

func New(cfg *config.Config) *OpsManager {
	return &OpsManager{cfg: cfg}
}

func (om *OpsManager) SetupAuth() error {
	unauthenticatedClient := network.NewUnauthenticatedClient(om.target(), config.SkipSSLValidation, time.Duration(requestTimeout)*time.Second)
	setupService := api.NewSetupService(unauthenticatedClient)

	cmd := commands.NewConfigureAuthentication(setupService, log.New(os.Stdout, "", 0))
	return cmd.Execute([]string{"--username", om.cfg.OpsManUsername, "--password", om.cfg.OpsManPassword, "--decryption-passphrase", om.cfg.OpsManDecryptionPhrase})
}

func (om *OpsManager) SetupBosh() error {
	client, err := om.oauthClient()
	if err != nil {
		return err
	}

	boshService := api.NewBoshFormService(client)
	diagnosticService := api.NewDiagnosticService(client)
	cmd := commands.NewConfigureBosh(boshService, diagnosticService, log.New(os.Stdout, "", 0))

	// TODO(jrjohnson): just hard coding this for testing
	return cmd.Execute([]string{"--iaas-configuration", "{\"project\": \"google.com:graphite-test-bosh-cpi-cert\",\"default_deployment_tag\": \"foo-vms\",\"auth_json\": \"\"}"})
}

// {"iaas_configuration":{"project":null,"associated_service_account":null},"director_configuration":{"ntp_servers_string":"","metrics_ip":null,"resurrector_enabled":false,"max_threads":null,"database_type":"internal","blobstore_type":"local"},"security_configuration":{"trusted_certificates":null,"generate_vm_passwords":true},"syslog_configuration":{"enabled":false}}

func (om *OpsManager) target() string {
	return fmt.Sprintf("https://%s", om.cfg.OpsManagerIp)
}

func (om *OpsManager) oauthClient() (network.OAuthClient, error) {
	return network.NewOAuthClient(om.target(), om.cfg.OpsManUsername, om.cfg.OpsManPassword, "", "", config.SkipSSLValidation, true, time.Duration(requestTimeout)*time.Second)
}
