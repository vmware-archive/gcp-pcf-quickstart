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
	"context"
	"log"

	"omg-cli/config"
	"omg-cli/templates"

	"github.com/alecthomas/kingpin"
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
	ctx := context.Background()
	envCfg, err := config.FromEnvDirectory(cmd.envDir)
	if err != nil {
		return err
	}

	envCfg.PivnetAPIToken = cmd.pivnetAPIToken

	mover, err := getMover(envCfg, cmd.tileCacheDir, cmd.logger)
	if err != nil {
		return err
	}

	pattern, err := templates.GetPattern(envCfg, map[string]interface{}{}, "", false)
	if err != nil {
		return err
	}

	pattern.Validate(false)
	if err != nil {
		return err
	}

	for _, tile := range pattern.Tiles {
		err = mover.Cache(ctx, tile.Product)
		if err != nil {
			return err
		}

		err = mover.Cache(ctx, tile.Stemcell)
		if err != nil {
			return err
		}
	}

	return nil
}
