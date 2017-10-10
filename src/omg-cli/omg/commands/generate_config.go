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

	"encoding/json"

	"fmt"
	"io/ioutil"

	"github.com/alecthomas/kingpin"
)

type GenerateConfigCommand struct {
	logger *log.Logger
	envDir string
}

const GenerateConfigCommandName = "generate-config"

func (cmd *GenerateConfigCommand) register(app *kingpin.Application) {
	c := app.Command(GenerateConfigCommandName, "Generate default environment configuration").Action(cmd.run)
	registerEnvConfigFlag(c, &cmd.envDir)
}

func (cmd *GenerateConfigCommand) run(c *kingpin.ParseContext) error {
	cfg, err := config.DefaultEnvConfig()
	if err != nil {
		return err
	}

	cfgStr, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	cfgPath := fmt.Sprintf("%s/config.json", cmd.envDir)
	if err := ioutil.WriteFile(cfgPath, cfgStr, 0644); err != nil {
		return err
	}

	cmd.logger.Printf("successfully wrote default config to: %s", cfgPath)

	return nil
}
