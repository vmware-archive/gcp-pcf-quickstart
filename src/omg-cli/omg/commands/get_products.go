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

	"encoding/json"
	"os"

	"github.com/alecthomas/kingpin"
)

type GetProductsCommand struct {
	logger              *log.Logger
	terraformConfigPath string
	opsManCreds         config.OpsManagerCredentials
}

const GetProductsName = "get-products"

func (dic *GetProductsCommand) register(app *kingpin.Application) {
	c := app.Command(GetProductsName, "List installed products").Action(dic.run)
	registerTerraformConfigFlag(c, &dic.terraformConfigPath)
	registerOpsManagerFlags(c, &dic.opsManCreds)
}

func (dic *GetProductsCommand) run(c *kingpin.ParseContext) error {
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

	dic.logger.Printf("fetched products")

	val, err := json.Marshal(&products)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "%s", string(val))

	fmt.Println()

	return nil
}
