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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"omg-cli/config"
	"omg-cli/templates"

	"github.com/alecthomas/kingpin"
	"github.com/iancoleman/strcase"
)

// SourceConfigCommand outputs the quickstart's config.
type SourceConfigCommand struct {
	logger *log.Logger
	envDir string
}

const sourceConfigCommandName = "source-config"

func (cmd *SourceConfigCommand) register(app *kingpin.Application) {
	c := app.Command(sourceConfigCommandName, "Output environment config as environment variables").Action(cmd.run)
	registerEnvConfigFlag(c, &cmd.envDir)
}

func nameToEnv(name string) string {
	return strcase.ToScreamingSnake(name)
}

func (cmd *SourceConfigCommand) run(c *kingpin.ParseContext) error {
	cfg, err := config.FromEnvDirectory(cmd.envDir)
	if err != nil {
		return err
	}

	opsmanImage, err := templates.Templates.Open("opsman-image")
	if err != nil {
		return fmt.Errorf("opening opsman image url: %v", err)
	}
	url, err := ioutil.ReadAll(opsmanImage)
	if err != nil {
		return fmt.Errorf("reading opsman image url: %v", err)
	}

	cfg.BaseImageURL = strings.TrimSpace(string(url))

	cfgBytes, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	flattened := map[string]interface{}{}
	if err = json.Unmarshal(cfgBytes, &flattened); err != nil {
		return err
	}

	for key, value := range flattened {
		cmd.logger.Printf(`%s="%v"`, nameToEnv(key), value)
	}

	return nil
}
