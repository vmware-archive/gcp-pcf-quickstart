package commands

import (
	"fmt"
	"log"
	"omg-cli/config"
	"omg-cli/omg/setup"
	"omg-cli/ops_manager"

	"github.com/alecthomas/kingpin"
)

type ConfigureOpsManagerCommand struct {
	logger              *log.Logger
	terraformConfigPath string
	opsManCreds         config.OpsManagerCredentials
}

func (comc *ConfigureOpsManagerCommand) register(app *kingpin.Application) {
	c := app.Command("ConfigureOpsManager", "Push desired tiles to a fresh Ops Manager for image capture").Action(comc.run)
	c.Flag("terraform-output-path", "JSON output from terraform state for deployment").Default("env.json").StringVar(&comc.terraformConfigPath)
	c.Flag("opsman-username", "Username for Ops Manager").Default(defaultUsername).StringVar(&comc.opsManCreds.Username)
	c.Flag("opsman-password", "Password for Ops Manager").Default(defaultPassword).StringVar(&comc.opsManCreds.Password)
	c.Flag("opsman-decryption-phrase", "Decryption Phrase for Ops Manager").Default(defaultDecryptionPhrase).StringVar(&comc.opsManCreds.DecryptionPhrase)
	c.Flag("opsman-skip-ssl-verification", "Skip SSL Validation for Ops Manager").Default(defaultSkipSSLVerify).BoolVar(&comc.opsManCreds.SkipSSLVerification)
}

func (comc *ConfigureOpsManagerCommand) run(c *kingpin.ParseContext) error {
	cfg, err := config.FromTerraform(comc.terraformConfigPath)
	if err != nil {
		return err
	}

	omSdk, err := ops_manager.NewSdk(fmt.Sprintf("https://%s", cfg.OpsManagerIp), comc.opsManCreds, *comc.logger)
	if err != nil {
		return err
	}

	opsMan := setup.NewService(cfg, omSdk, nil, comc.logger, selectedTiles)

	return runSteps([]step{
		opsMan.PoolTillOnline,
		opsMan.Unlock,
		opsMan.SetupBosh,
		opsMan.ConfigureTiles,
		opsMan.ApplyChanges,
	})
}
