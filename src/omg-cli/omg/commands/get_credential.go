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
	"omg-cli/ops_manager"

	"github.com/alecthomas/kingpin"
)

// GetCredentialCommand prints credentials.
type GetCredentialCommand struct {
	logger     *log.Logger
	envDir     string
	appName    string
	credential string
}

const getCredentialName = "get-credential"

func (cmd *GetCredentialCommand) register(app *kingpin.Application) {
	c := app.Command(getCredentialName, "Fetch a credential for a tile").Action(cmd.run)
	registerEnvConfigFlag(c, &cmd.envDir)

	c.Flag("app-name", "Name of the Product (type)").Required().StringVar(&cmd.appName)
	c.Flag("credential", "Credential to fetch (eg .uaa.admin_credentials)").Required().StringVar(&cmd.credential)
}

func (cmd *GetCredentialCommand) run(c *kingpin.ParseContext) error {
	cfg, err := config.TerraformFromEnvDirectory(cmd.envDir)
	if err != nil {
		return err
	}

	omSdk, err := ops_manager.NewSdk(fmt.Sprintf("https://%s", cfg.OpsManagerHostname), cfg.OpsManager, cmd.logger)
	if err != nil {
		return err
	}

	cred, err := omSdk.GetCredentials(cmd.appName, cmd.credential)
	if err != nil {
		return err
	}

	cmd.logger.Printf("identity: %s", cred.Identity)
	cmd.logger.Printf("password: %s", cred.Password)

	return nil
}
