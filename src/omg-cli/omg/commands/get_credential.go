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

type GetCredentialCommand struct {
	logger              *log.Logger
	terraformConfigPath string
	opsManCreds         config.OpsManagerCredentials
	productType         string
	credential          string
}

const GetCredentialName = "get-credential"

func (dic *GetCredentialCommand) register(app *kingpin.Application) {
	c := app.Command(GetCredentialName, "Fetch a credential for a tile").Action(dic.run)
	registerTerraformConfigFlag(c, &dic.terraformConfigPath)
	registerOpsManagerFlags(c, &dic.opsManCreds)

	c.Flag("app-name", "Name of the Product (type)").Required().StringVar(&dic.productType)
	c.Flag("credential", "Credential to fetch (eg .uaa.admin_credentials)").Required().StringVar(&dic.credential)
}

func (dic *GetCredentialCommand) run(c *kingpin.ParseContext) error {
	cfg, err := config.FromTerraform(dic.terraformConfigPath)
	if err != nil {
		return err
	}

	omSdk, err := ops_manager.NewSdk(fmt.Sprintf("https://%s", cfg.OpsManagerIp), dic.opsManCreds, *dic.logger)
	if err != nil {
		return err
	}

	products, err := omSdk.GetProducts()
	if err != nil {
		return err
	}

	appGuid := ""
	for _, p := range products {
		if p.Type == dic.productType {
			appGuid = p.Guid
		}
	}

	if appGuid == "" {
		return fmt.Errorf("could not find installed application by name: %s", dic.productType)
	}

	cred, err := omSdk.GetCredentials(appGuid, dic.credential)
	if err != nil {
		return err
	}

	dic.logger.Printf("identity: %s", cred.Identity)
	dic.logger.Printf("password: %s", cred.Password)

	return nil
}
