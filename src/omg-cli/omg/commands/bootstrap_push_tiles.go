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

type BootstrapPushTilesCommand struct {
	logger              *log.Logger
	apiToken            string
	terraformConfigPath string
	username            string
	sshKeyPath          string
}

const BootstrapPushTilesName = "bootstrap-push-tiles"

func (bj *BootstrapPushTilesCommand) register(app *kingpin.Application) {
	c := app.Command(BootstrapPushTilesName, "Prepare Ops Manager for image capture from outside the network").Action(bj.run)
	c.Flag("username", "Username to login on jumpbox").Required().StringVar(&bj.username)
	c.Flag("ssh-key-path", "Path to SSH to login on jumpbox").Required().StringVar(&bj.sshKeyPath)
	registerTerraformConfigFlag(c, &bj.terraformConfigPath)
	registerPivnetFlag(c, &bj.apiToken)
}

func (bj *BootstrapPushTilesCommand) run(c *kingpin.ParseContext) error {
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
			return jb.RunOmg(fmt.Sprintf("%s --pivnet-api-token=%s --terraform-output-path=env.json", PushTilesName, bj.apiToken))
		},
	})
}
