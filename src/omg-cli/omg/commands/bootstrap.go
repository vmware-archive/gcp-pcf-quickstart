package commands

import (
	"log"
	"omg-cli/config"

	"omg-cli/omg/setup"

	"fmt"

	"github.com/alecthomas/kingpin"
)

type BootstrapCommand struct {
	logger              *log.Logger
	apiToken            string
	terraformConfigPath string
	username            string
	sshKeyPath          string
}

func (bj *BootstrapCommand) register(app *kingpin.Application) {
	c := app.Command("bootstrap", "Deploy OMG on provisioned infrastructure from outside the network").Action(bj.run)
	c.Flag("username", "Username to login on jumpbox").Required().StringVar(&bj.username)
	c.Flag("ssh-key-path", "Path to SSH to login on jumpbox").Required().StringVar(&bj.sshKeyPath)
	registerTerraformConfigFlag(c, &bj.terraformConfigPath)
}

func (bj *BootstrapCommand) run(c *kingpin.ParseContext) error {
	cfg, err := config.FromTerraform(bj.terraformConfigPath)
	if err != nil {
		return fmt.Errorf("load terraform config: %v", err)
	}

	jb, err := setup.NewJumpbox(bj.logger, cfg.JumpboxIp, bj.username, bj.sshKeyPath, bj.terraformConfigPath)
	if err != nil {
		return fmt.Errorf("connect to jumpbox: %v", err)
	}

	return runSteps([]step{
		jb.PoolTillStarted,
		jb.UploadDependencies,
		jb.ConfigureOpsManager,
	})
}
