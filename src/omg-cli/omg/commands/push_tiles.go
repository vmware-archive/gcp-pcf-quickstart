package commands

import (
	"fmt"
	"log"

	"omg-cli/config"
	"omg-cli/omg/setup"
	"omg-cli/ops_manager"
	"omg-cli/pivnet"

	"github.com/alecthomas/kingpin"
)

type PushTilesCommand struct {
	logger              *log.Logger
	apiToken            string
	terraformConfigPath string
	opsManCreds         config.OpsManagerCredentials
}

const PushTilesName = "push-tiles"

func (bic *PushTilesCommand) register(app *kingpin.Application) {
	c := app.Command(PushTilesName, "Push desired tiles to a deployed Ops Manager").Action(bic.run)
	registerTerraformConfigFlag(c, &bic.terraformConfigPath)
	registerOpsManagerFlags(c, &bic.opsManCreds)
	registerPivnetFlag(c, &bic.apiToken)
}

func (bic *PushTilesCommand) run(c *kingpin.ParseContext) error {
	cfg, err := config.FromTerraform(bic.terraformConfigPath)
	if err != nil {
		return err
	}

	omSdk, err := ops_manager.NewSdk(fmt.Sprintf("https://%s", cfg.OpsManagerIp), bic.opsManCreds, *bic.logger)
	if err != nil {
		return err
	}

	pivnetSdk, err := pivnet.NewSdk(bic.apiToken, bic.logger)
	if err != nil {
		return err
	}

	opsMan := setup.NewService(cfg, omSdk, pivnetSdk, bic.logger, selectedTiles)

	return runSteps([]step{
		opsMan.PoolTillOnline,
		opsMan.SetupAuth,
		opsMan.UploadTiles,
	})
}
