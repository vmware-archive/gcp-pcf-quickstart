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
	"os"
	"path/filepath"

	"omg-cli/config"
	"omg-cli/version"

	"github.com/alecthomas/kingpin"

	"github.com/starkandwayne/om-tiler/mover"
	"github.com/starkandwayne/om-tiler/opsman"
	"github.com/starkandwayne/om-tiler/pivnet"
	"github.com/starkandwayne/om-tiler/tiler"
)

type register interface {
	register(app *kingpin.Application)
}

// Configure sets up the kingpin commands for the omg-cli.
func Configure(logger *log.Logger, app *kingpin.Application) {
	cmds := []register{
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

func getPivnet(envCfg *config.EnvConfig, l *log.Logger) *pivnet.Client {
	return pivnet.NewClient(pivnet.Config{
		Token: envCfg.PivnetAPIToken,
		UserAgent:  version.UserAgent(),
		AcceptEULA: true,
	}, l)
}

func getMover(envCfg *config.EnvConfig, c string, l *log.Logger) (*mover.Mover, error) {
	if _, err := os.Stat(c); os.IsNotExist(err) {
		if err := os.Mkdir(c, os.ModePerm); err != nil {
			return nil, fmt.Errorf("creating tile cache directory %s: %v", c, err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("finding tile cache directory %s: %v", c, err)
	}

	return mover.NewMover(getPivnet(envCfg, l), c, l)
}

func getTiler(cfg *config.Config, envCfg *config.EnvConfig, c string, l *log.Logger) (*tiler.Tiler, error) {
	omClient, err := opsman.NewClient(opsman.Config{
		Target:               cfg.OpsManagerHostname,
		Username:             cfg.OpsManager.Username,
		Password:             cfg.OpsManager.Password,
		DecryptionPassphrase: cfg.OpsManager.DecryptionPhrase,
		SkipSSLVerification:  cfg.OpsManager.SkipSSLVerification,
	}, l)
	if err != nil {
		return nil, err
	}

	mover, err := getMover(envCfg, c, l)
	if err != nil {
		return nil, err
	}

	return tiler.NewTiler(omClient, mover, l), nil
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
