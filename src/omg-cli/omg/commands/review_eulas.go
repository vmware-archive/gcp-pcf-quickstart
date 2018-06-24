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
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"omg-cli/config"
	"omg-cli/pivnet"

	"github.com/alecthomas/kingpin"
)

type ReviewEulasCommand struct {
	logger    *log.Logger
	envDir    string
	envConfig *config.EnvConfig
	acceptAll bool
	pivnetSdk *pivnet.Sdk
}

const ReviewEulasName = "review-eulas"

var eulaSlugs = []string{"pivotal_software_eula"}

func (cmd *ReviewEulasCommand) register(app *kingpin.Application) {
	c := app.Command(ReviewEulasName, "View product EULAs and interactively accept/deny").Action(cmd.run)
	registerEnvConfigFlag(c, &cmd.envDir)
	c.Flag("accept-all", "Accept all EULAs non-interactively").Default("false").BoolVar(&cmd.acceptAll)
}

func (cmd *ReviewEulasCommand) run(c *kingpin.ParseContext) error {
	var err error
	cmd.envConfig, err = config.ConfigFromEnvDirectory(cmd.envDir)
	if err != nil {
		cmd.logger.Fatalf("loading environment config: %v", err)
	}

	cmd.pivnetSdk, err = pivnet.NewSdk(cmd.envConfig.PivnetApiToken, cmd.logger)
	if err != nil {
		return err
	}

	return run([]step{
		cmd.fetchAndPrompt,
		cmd.acceptEulas,
	})
}

func (cmd *ReviewEulasCommand) fetchAndPrompt() error {
	var eulas []*pivnet.Eula
	for _, slug := range eulaSlugs {
		eula, err := cmd.pivnetSdk.GetEula(slug)
		if err != nil {
			return err
		}

		eulas = append(eulas, eula)
	}

	reader := bufio.NewReader(os.Stdin)
	for _, eula := range eulas {
		fmt.Printf("EULA: %s\n%s\n", eula.Name, eula.Content)

		if cmd.acceptAll {
			fmt.Printf("EULA accepted via command line flag\n")
		} else {
			fmt.Printf("Accept EULA? (y/n): ")

			input, _ := reader.ReadString('\n')
			if !strings.HasPrefix(strings.ToLower(input), "y") {
				return fmt.Errorf("can not proceed without EULA concent")
			}
		}
	}

	return nil
}

func (cmd *ReviewEulasCommand) acceptEulas() error {
	var tileData []config.PivnetMetadata
	tiles := selectedTiles(cmd.logger, cmd.envConfig)
	for _, installer := range tiles {
		if installer.BuiltIn() {
			continue
		}

		tile := installer.Definition(cmd.envConfig)

		tileData = append(tileData, tile.Pivnet)
		if tile.Stemcell != nil {
			tileData = append(tileData, tile.Stemcell.PivnetMetadata)
		}
	}

	for _, tile := range tileData {
		if err := cmd.pivnetSdk.AcceptEula(tile); err != nil {
			return fmt.Errorf("accepting EULA for %s: %v", tile.Name, err)
		}
	}

	cmd.logger.Printf("accepted EULAs for %d products", len(tileData))

	return nil
}
