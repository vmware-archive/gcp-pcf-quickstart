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
	"log"

	"omg-cli/config"
	"omg-cli/templates"

	"omg-cli/version"

	"github.com/alecthomas/kingpin"

	"github.com/starkandwayne/om-tiler/mover"
	"github.com/starkandwayne/om-tiler/opsman"
	"github.com/starkandwayne/om-tiler/pattern"
	"github.com/starkandwayne/om-tiler/pivnet"
	"github.com/starkandwayne/om-tiler/tiler"
)

// PushTilesCommand pushes tiles to the Ops Manager.
type PushTilesCommand struct {
	logger       *log.Logger
	envDir       string
	tileCacheDir string
}

const pushTilesName = "push-tiles"

func (cmd *PushTilesCommand) register(app *kingpin.Application) {
	c := app.Command(pushTilesName, "Push desired tiles to a deployed Ops Manager").Action(cmd.run)
	registerEnvConfigFlag(c, &cmd.envDir)
	registerTileCacheFlag(c, &cmd.tileCacheDir)
}

func (cmd *PushTilesCommand) run(c *kingpin.ParseContext) error {
	cfg, err := config.TerraformFromEnvDirectory(cmd.envDir)
	if err != nil {
		return err
	}

	envCfg, err := config.FromEnvDirectory(cmd.envDir)
	if err != nil {
		return err
	}

	omclient, err := opsman.NewClient(opsman.Config{
		Target:               cfg.OpsManagerHostname,
		Username:             cfg.OpsManager.Username,
		Password:             cfg.OpsManager.Password,
		DecryptionPassphrase: cfg.OpsManager.DecryptionPhrase,
		SkipSSLVerification:  cfg.OpsManager.SkipSSLVerification,
	}, cmd.logger)
	if err != nil {
		return err
	}

	pivnetclient := pivnet.NewClient(pivnet.Config{
		Token:      envCfg.PivnetAPIToken,
		UserAgent:  version.UserAgent(),
		AcceptEULA: true,
	}, cmd.logger)

	mover, err := mover.NewMover(pivnetclient, "", cmd.logger)
	if err != nil {
		return err
	}

	omConfigurator, err := tiler.NewTiler(omclient, mover, cmd.logger)
	if err != nil {
		return err
	}
	var opsFiles []string
	if envCfg.SmallFootprint {
		opsFiles = append(opsFiles, "options/small-footprint.yml")
	}
	if envCfg.IncludeHealthwatch {
		opsFiles = append(opsFiles, "options/healthwatch.yml")
	}
	return omConfigurator.Apply(pattern.Template{
		Store:    templates.Templates,
		Manifest: "deployment.yml",
		OpsFiles: opsFiles,
		Vars:     cfg.Raw,
	})

}
