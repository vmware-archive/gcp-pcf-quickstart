package commands

import (
	"fmt"
	"log"

	"omg-cli/config"
	"omg-cli/omg/setup"

	"github.com/alecthomas/kingpin"
)

type BootstrapDeployCommand struct {
	logger              *log.Logger
	terraformConfigPath string
	username            string
	sshKeyPath          string
}

const BootstrapDeployCommandName = "bootstrap-deploy"

func (bj *BootstrapDeployCommand) register(app *kingpin.Application) {
	c := app.Command(BootstrapDeployCommandName, "Deploy PCF on provisioned infrastructure from outside the network").Action(bj.run)
	c.Flag("username", "Username to login on jumpbox").Required().StringVar(&bj.username)
	c.Flag("ssh-key-path", "Path to SSH to login on jumpbox").Required().StringVar(&bj.sshKeyPath)
	registerTerraformConfigFlag(c, &bj.terraformConfigPath)
}

func (bj *BootstrapDeployCommand) run(c *kingpin.ParseContext) error {
	cfg, err := config.FromTerraform(bj.terraformConfigPath)
	if err != nil {
		return fmt.Errorf("load terraform config: %v", err)
	}

	jb, err := setup.NewJumpbox(bj.logger, cfg.JumpboxIp, bj.username, bj.sshKeyPath, bj.terraformConfigPath)
	if err != nil {
		return fmt.Errorf("connect to jumpbox: %v", err)
	}

	return run([]step{
		jb.PoolTillStarted,
		jb.UploadDependencies,
		func() error {
			return jb.RunOmg(fmt.Sprintf("%s --terraform-output-path=env.json", DeployName))
		},
	})
}
