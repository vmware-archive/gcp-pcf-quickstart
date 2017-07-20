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

const ConfigureOpsManagerCommandName = "configure"

func (comc *ConfigureOpsManagerCommand) register(app *kingpin.Application) {
	c := app.Command(ConfigureOpsManagerCommandName, "Push desired tiles to a fresh Ops Manager for image capture").Action(comc.run)
	registerTerraformConfigFlag(c, &comc.terraformConfigPath)
	registerOpsManagerFlags(c, &comc.opsManCreds)
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
