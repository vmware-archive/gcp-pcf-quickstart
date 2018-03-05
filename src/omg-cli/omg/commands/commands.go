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
	"path/filepath"

	"omg-cli/omg/tiles"
	"omg-cli/omg/tiles/ert"
	"omg-cli/omg/tiles/gcp_director"
	"omg-cli/omg/tiles/service_broker"
	"omg-cli/omg/tiles/stackdriver_nozzle"

	"fmt"

	"github.com/alecthomas/kingpin"
)

const (
	defaultSkipSSLVerify = "true"
)

// TODO(jrjohnson): Remove? Move?
var selectedTiles []tiles.TileInstaller

type register interface {
	register(app *kingpin.Application)
}

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

	selectedTiles = []tiles.TileInstaller{
		&gcp_director.Tile{},
		&ert.Tile{},
		&stackdriver_nozzle.Tile{Logger: logger},
		&service_broker.Tile{},
		// TODO: enable conditionally
		//&healthwatch.Tile{},
	}

}

type step func() error

func run(steps []step) error {
	for _, v := range steps {
		if err := v(); err != nil {
			return err
		}
	}
	return nil
}

func retry(fn step, times int) error {
	errs := []error{}

	for i := 0; i < times; i++ {
		if err := fn(); err != nil {
			errs = append(errs, err)
		} else {
			return nil
		}
	}

	return fmt.Errorf("failed after %d attempts, errors: %v", times, errs)
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
