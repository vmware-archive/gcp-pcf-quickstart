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
	"path/filepath"

	"omg-cli/config"
	"omg-cli/omg/tiles"
	"omg-cli/omg/tiles/ert"
	"omg-cli/omg/tiles/gcp_director"
	"omg-cli/omg/tiles/healthwatch"
	"omg-cli/omg/tiles/service_broker"
	"omg-cli/omg/tiles/stackdriver_nozzle"

	"github.com/alecthomas/kingpin"
)

type register interface {
	register(app *kingpin.Application)
}

// Configure sets up the kingpin commands for the omg-cli.
func Configure(logger *log.Logger, app *kingpin.Application) {
	cmds := []register{
		&PushTilesCommand{logger: logger},
		&DeployCommand{logger: logger},
		&DeleteInstallationCommand{logger: logger},
		&GetCredentialCommand{logger: logger},
		&RemoteCommand{logger: logger},
		&PrepareProjectCommand{logger: logger},
		&GenerateConfigCommand{logger: logger},
		&SourceConfigCommand{logger: logger},
		&ReviewEulasCommand{logger: logger},
		&CleanupProjectCommand{logger: logger},
		&DirectorSSHCommand{logger: logger},
		&CacheTilesCommand{logger: logger},
	}

	for _, c := range cmds {
		c.register(app)
	}
}

func selectedTiles(logger *log.Logger, config *config.EnvConfig) []tiles.TileInstaller {
	result := []tiles.TileInstaller{
		&gcp_director.Tile{},
		&ert.Tile{},
		&stackdriver_nozzle.Tile{Logger: logger},
		&service_broker.Tile{},
	}
	if config.IncludeHealthwatch {
		result = append(result, &healthwatch.Tile{})
	}
	return result
}

type step struct {
	function func() error
	name     string
}

func run(steps []step, logger *log.Logger) error {
	for _, v := range steps {
		logger.Printf("running step: %s", v.name)
		if err := v.function(); err != nil {
			return fmt.Errorf("failed running step %s, error: %v", v.name, err)
		}
	}
	return nil
}

func registerEnvConfigFlag(c *kingpin.CmdClause, path *string) {
	c.Flag("env-dir", "path to environment configuration and state").Default(filepath.Join("env", "pcf")).StringVar(path)
}

func registerTileCacheFlag(c *kingpin.CmdClause, path *string) {
	c.Flag("cache-dir", "path to directory used to cache downloads").Default("cache").StringVar(path)
}
func registerQuietFlag(c *kingpin.CmdClause, quiet *bool) {
	c.Flag("quiet", "quiet output, no non-essential information").Default("false").BoolVar(quiet)
}

func registerPivnetAPITokenFlag(c *kingpin.CmdClause, token *string) {
	c.Flag("pivnet-api-token", "API token for network.pivotal.io (see: https://network.pivotal.io/users/dashboard/edit-profile)").Required().StringVar(token)
}
