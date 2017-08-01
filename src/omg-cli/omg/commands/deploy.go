package commands

import (
	"fmt"
	"log"

	"omg-cli/config"
	"omg-cli/omg/setup"
	"omg-cli/ops_manager"

	"github.com/alecthomas/kingpin"
)

type DeployCommand struct {
	logger              *log.Logger
	terraformConfigPath string
	opsManCreds         config.OpsManagerCredentials
	applyChanges        bool
}

const DeployName = "deploy"

func (dc *DeployCommand) register(app *kingpin.Application) {
	c := app.Command(DeployName, "Deploy tiles to a freshly deployed Ops Manager").Action(dc.run)
	registerTerraformConfigFlag(c, &dc.terraformConfigPath)
	registerOpsManagerFlags(c, &dc.opsManCreds)
	c.Flag("apply-changes", "Apply Changes").Default("true").BoolVar(&dc.applyChanges)
}

func (dc *DeployCommand) run(c *kingpin.ParseContext) error {
	cfg, err := config.FromTerraform(dc.terraformConfigPath)
	if err != nil {
		return err
	}

	omSdk, err := ops_manager.NewSdk(fmt.Sprintf("https://%s", cfg.OpsManagerIp), dc.opsManCreds, *dc.logger)
	if err != nil {
		return err
	}

	opsMan := setup.NewService(cfg, omSdk, nil, dc.logger, selectedTiles)

	steps := []step{
		opsMan.PoolTillOnline,
		opsMan.Unlock,
		opsMan.ConfigureTiles,
	}

	if dc.applyChanges {
		steps = append(steps, opsMan.ApplyChanges)
	}

	return run(steps)
}
