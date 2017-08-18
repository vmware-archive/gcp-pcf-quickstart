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

	omSdk, err := ops_manager.NewSdk(fmt.Sprintf("https://%s", cfg.OpsManagerHostname), dc.opsManCreds, *dc.logger)
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
