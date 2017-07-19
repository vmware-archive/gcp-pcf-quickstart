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

type BakeImageCommand struct {
	logger              *log.Logger
	apiToken            string
	terraformConfigPath string
	opsManCreds         config.OpsManagerCredentials
}

func (bic *BakeImageCommand) register(app *kingpin.Application) {
	c := app.Command("BakeImage", "Push desired tiles to a fresh Ops Manager for image capture").Action(bic.run)
	c.Flag("pivnet-api-token", "Look for 'API TOKEN' at https://network.pivotal.io/users/dashboard/edit-profile.").Required().StringVar(&bic.apiToken)
	registerTerraformConfigFlag(c, &bic.terraformConfigPath)
	registerOpsManagerFlags(c, &bic.opsManCreds)
}

func (bic *BakeImageCommand) run(c *kingpin.ParseContext) error {
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
