/*
 * Copyright 2017 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
	applyChanges        bool
}

const BootstrapDeployCommandName = "bootstrap-deploy"

func (bdc *BootstrapDeployCommand) register(app *kingpin.Application) {
	c := app.Command(BootstrapDeployCommandName, "Deploy PCF on provisioned infrastructure from outside the network").Action(bdc.run)
	c.Flag("username", "Username to login on jumpbox").Required().StringVar(&bdc.username)
	c.Flag("ssh-key-path", "Path to SSH to login on jumpbox").Required().StringVar(&bdc.sshKeyPath)
	c.Flag("apply-changes", "Apply Changes").Default("true").BoolVar(&bdc.applyChanges)
	registerTerraformConfigFlag(c, &bdc.terraformConfigPath)
}

func (bdc *BootstrapDeployCommand) run(c *kingpin.ParseContext) error {
	cfg, err := config.FromTerraform(bdc.terraformConfigPath)
	if err != nil {
		return fmt.Errorf("load terraform config: %v", err)
	}

	jb, err := setup.NewJumpbox(bdc.logger, cfg.JumpboxIp, bdc.username, bdc.sshKeyPath, bdc.terraformConfigPath)
	if err != nil {
		return fmt.Errorf("connect to jumpbox: %v", err)
	}

	cmd := fmt.Sprintf("%s --terraform-output-path=env.json", DeployName)
	if !bdc.applyChanges {
		cmd += " --no-apply-changes"
	}

	return run([]step{
		jb.PoolTillStarted,
		jb.UploadDependencies,
		func() error {
			return jb.RunOmg(cmd)
		},
	})
}
