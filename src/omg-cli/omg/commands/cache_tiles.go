/*
 * Copyright 2018 Google Inc.
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
	"os"

	"omg-cli/config"
	"omg-cli/templates"
	"omg-cli/version"

	"github.com/alecthomas/kingpin"
	"github.com/starkandwayne/om-tiler/mover"
	"github.com/starkandwayne/om-tiler/pivnet"
)

// CacheTilesCommand caches tiles to the given tileCacheDir.
type CacheTilesCommand struct {
	logger         *log.Logger
	envDir         string
	tileCacheDir   string
	pivnetAPIToken string
}

const cacheTilesName = "cache-tiles"

func (cmd *CacheTilesCommand) register(app *kingpin.Application) {
	c := app.Command(cacheTilesName, "Cache tile downloads locally").Action(cmd.run)
	registerEnvConfigFlag(c, &cmd.envDir)
	registerTileCacheFlag(c, &cmd.tileCacheDir)
	registerPivnetAPITokenFlag(c, &cmd.pivnetAPIToken)
}

func (cmd *CacheTilesCommand) run(c *kingpin.ParseContext) error {
	envCfg, err := config.FromEnvDirectory(cmd.envDir)
	if err != nil {
		return err
	}

	pivnetClient := pivnet.NewClient(pivnet.Config{
		Token:      cmd.pivnetAPIToken,
		UserAgent:  version.UserAgent(),
		AcceptEULA: true,
	}, cmd.logger)

	if _, err := os.Stat(cmd.tileCacheDir); os.IsNotExist(err) {
		if err := os.Mkdir(cmd.tileCacheDir, os.ModePerm); err != nil {
			return fmt.Errorf("creating tile cache directory %s: %v", cmd.tileCacheDir, err)
		}
	} else if err != nil {
		return fmt.Errorf("finding tile cache directory %s: %v", cmd.tileCacheDir, err)
	}

	mover, err := mover.NewMover(pivnetClient, cmd.tileCacheDir, cmd.logger)
	if err != nil {
		return err
	}

	pattern, err := templates.GetPattern(envCfg, map[string]interface{}{})
	if err != nil {
		return err
	}

	pattern.Validate(false)
	if err != nil {
		return err
	}

	for _, tile := range pattern.Tiles {
		err = mover.Cache(tile.Product)
		if err != nil {
			return err
		}

		err = mover.Cache(tile.Stemcell)
		if err != nil {
			return err
		}
	}

	return nil
}
