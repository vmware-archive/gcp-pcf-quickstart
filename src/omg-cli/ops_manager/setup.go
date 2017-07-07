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

func (om *OpsManager) SetupAuth() {
	unauthenticatedClient := network.NewUnauthenticatedClient(om.target(), config.SkipSSLValidation, time.Duration(requestTimeout)*time.Second)
	setupService := api.NewSetupService(unauthenticatedClient)

	stdout := log.New(os.Stdout, "", 0)

	cmd := commands.NewConfigureAuthentication(setupService, stdout)
	cmd.Execute([]string{"--username", cfg.OpsManUsername, "--password", cfg.OpsManPassword, "--decryption-passphrase", cfg.OpsManDecryptionPhrase})
}

func (om *OpsManager) target() string {
	return fmt.Sprintf("https://%s", om.cfg.OpsManagerIp)
}
