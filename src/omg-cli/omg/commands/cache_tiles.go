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
	"path/filepath"

	"omg-cli/config"
	"omg-cli/pivnet"

	"github.com/alecthomas/kingpin"
)

type CacheTilesCommand struct {
	logger         *log.Logger
	envDir         string
	tileCacheDir   string
	pivnetApiToken string
}

const CacheTilesName = "cache-tiles"

func (cmd *CacheTilesCommand) register(app *kingpin.Application) {
	c := app.Command(CacheTilesName, "Cache tile downloads locally").Action(cmd.run)
	registerEnvConfigFlag(c, &cmd.envDir)
	registerTileCacheFlag(c, &cmd.tileCacheDir)
	registerPivnetApiTokenFlag(c, &cmd.pivnetApiToken)
}

func (cmd *CacheTilesCommand) run(c *kingpin.ParseContext) error {
	pivnetSdk, err := pivnet.NewSdk(cmd.pivnetApiToken, cmd.logger)
	if err != nil {
		return err
	}

	envCfg, err := config.ConfigFromEnvDirectory(cmd.envDir)
	if err != nil {
		return err
	}

	if _, err := os.Stat(cmd.tileCacheDir); os.IsNotExist(err) {
		if err := os.Mkdir(cmd.tileCacheDir, os.ModePerm); err != nil {
			return fmt.Errorf("creating tile cache directory %s: %v", cmd.tileCacheDir, err)
		}
	} else if err != nil {
		return fmt.Errorf("finding tile cache directory %s: %v", cmd.tileCacheDir, err)
	}

	tileCache := pivnet.TileCache{cmd.tileCacheDir}
	tiles := selectedTiles(cmd.logger, envCfg)
	for _, tile := range tiles {
		if tile.BuiltIn() {
			continue
		}
		definition := tile.Definition(&config.EnvConfig{SmallFootprint: true})
		cmd.logger.Printf("caching tile: %s", definition.Product.Name)

		output := filepath.Join(cmd.tileCacheDir, tileCache.FileName(definition.Pivnet))
		file, err := pivnetSdk.DownloadTileToPath(definition.Pivnet, output)
		if err != nil {
			return fmt.Errorf("downloading tile: %v", err)
		}
		if err := file.Close(); err != nil {
			return fmt.Errorf("closing tile: %v", err)
		}
	}

	return nil
}
