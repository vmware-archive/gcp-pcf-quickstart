package commands

import (
	"fmt"
	"log"

	"omg-cli/config"
	"omg-cli/omg/setup"
	"omg-cli/ops_manager"

	"github.com/alecthomas/kingpin"
)

type Deploy struct {
	logger              *log.Logger
	terraformConfigPath string
	opsManCreds         config.OpsManagerCredentials
}

const DeployName = "deploy"

func (comc *Deploy) register(app *kingpin.Application) {
	c := app.Command(DeployName, "Deploy tiles to a freshly deployed Ops Manager").Action(comc.run)
	registerTerraformConfigFlag(c, &comc.terraformConfigPath)
	registerOpsManagerFlags(c, &comc.opsManCreds)
}

func (comc *Deploy) run(c *kingpin.ParseContext) error {
	cfg, err := config.FromTerraform(comc.terraformConfigPath)
	if err != nil {
		return err
	}

	omSdk, err := ops_manager.NewSdk(fmt.Sprintf("https://%s", cfg.OpsManagerIp), comc.opsManCreds, *comc.logger)
	if err != nil {
		return err
	}

	opsMan := setup.NewService(cfg, omSdk, nil, comc.logger, selectedTiles)

	return run([]step{
		opsMan.PoolTillOnline,
		opsMan.Unlock,
		opsMan.ConfigureTiles,
		func() error { return retry(opsMan.ApplyChanges, 5) },
	})
}
