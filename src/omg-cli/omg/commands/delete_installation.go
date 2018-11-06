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
	"omg-cli/ops_manager"

	"github.com/alecthomas/kingpin"
)

type DeleteInstallationCommand struct {
	logger *log.Logger
	envDir string
}

const DeleteInstallationName = "delete-installation"

func (cmd *DeleteInstallationCommand) register(app *kingpin.Application) {
	c := app.Command(DeleteInstallationName, "Delete an Ops Manager installation").Action(cmd.run)
	registerEnvConfigFlag(c, &cmd.envDir)
}

func (cmd *DeleteInstallationCommand) run(c *kingpin.ParseContext) error {
	cfg, err := config.TerraformFromEnvDirectory(cmd.envDir)
	if err != nil {
		return err
	}

	envCfg, err := config.ConfigFromEnvDirectory(cmd.envDir)
	if err != nil {
		return err
	}

	omSdk, err := ops_manager.NewSdk(fmt.Sprintf("https://%s", cfg.OpsManagerHostname), cfg.OpsManager, *cmd.logger)
	if err != nil {
		return err
	}

	opsMan := setup.NewService(cfg, envCfg, omSdk, nil, cmd.logger, nil, nil)

	steps := []step{
		{opsMan.PoolTillOnline, "PoolTillOnline"},
		{opsMan.Unlock, "Unlock"},
		{opsMan.DeleteInstallation, "DeleteInstallation"},
	}

	return run(steps)
}
