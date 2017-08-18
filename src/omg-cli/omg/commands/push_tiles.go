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
	"omg-cli/pivnet"

	"github.com/alecthomas/kingpin"
)

type PushTilesCommand struct {
	logger              *log.Logger
	terraformConfigPath string
	opsManCreds         config.OpsManagerCredentials
}

const PushTilesName = "push-tiles"

func (bic *PushTilesCommand) register(app *kingpin.Application) {
	c := app.Command(PushTilesName, "Push desired tiles to a deployed Ops Manager").Action(bic.run)
	registerTerraformConfigFlag(c, &bic.terraformConfigPath)
	registerOpsManagerFlags(c, &bic.opsManCreds)
}

func (bic *PushTilesCommand) run(c *kingpin.ParseContext) error {
	cfg, err := config.FromTerraform(bic.terraformConfigPath)
	if err != nil {
		return err
	}

	omSdk, err := ops_manager.NewSdk(fmt.Sprintf("https://%s", cfg.OpsManagerHostname), bic.opsManCreds, *bic.logger)
	if err != nil {
		return err
	}

	pivnetSdk, err := pivnet.NewSdk(cfg.PivnetApiToken, bic.logger)
	if err != nil {
		return err
	}

	opsMan := setup.NewService(cfg, omSdk, pivnetSdk, bic.logger, selectedTiles)

	return run([]step{
		opsMan.PoolTillOnline,
		opsMan.SetupAuth,
		opsMan.UploadTiles,
	})
}
