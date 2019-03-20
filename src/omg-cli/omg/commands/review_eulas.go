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
	"omg-cli/templates"

	"github.com/alecthomas/kingpin"

	"github.com/starkandwayne/om-tiler/pattern"
	"github.com/starkandwayne/om-tiler/pivnet"
)

// ReviewEulasCommand reviews and accepts Pivotal's EULAs.
type ReviewEulasCommand struct {
	logger    *log.Logger
	envDir    string
	envConfig *config.EnvConfig
	acceptAll bool
	pivnet    *pivnet.Client
	pattern   pattern.Pattern
}

const reviewEulasName = "review-eulas"

func (cmd *ReviewEulasCommand) register(app *kingpin.Application) {
	c := app.Command(reviewEulasName, "View product EULAs and interactively accept/deny").Action(cmd.run)
	registerEnvConfigFlag(c, &cmd.envDir)
	c.Flag("accept-all", "Accept all EULAs non-interactively").Default("false").BoolVar(&cmd.acceptAll)
}

func (cmd *ReviewEulasCommand) run(c *kingpin.ParseContext) error {
	var err error
	cmd.envConfig, err = config.FromEnvDirectory(cmd.envDir)
	if err != nil {
		cmd.logger.Fatalf("loading environment config: %v", err)
	}

	cmd.pivnet = getPivnet(cmd.envConfig, cmd.logger)

	cmd.pattern, err = templates.GetPattern(cmd.envConfig, map[string]interface{}{})
	if err != nil {
		return err
	}

	return run([]step{
		{function: cmd.fetchAndPrompt, name: "fetchAndPrompt"},
		{function: cmd.acceptEulas, name: "acceptEulas"},
	}, cmd.logger)
	return nil
}

func (cmd *ReviewEulasCommand) fetchAndPrompt() error {
	eulas := make(map[string]*pivnet.EULA)
	for _, tile := range cmd.pattern.Tiles {
		eula, err := cmd.pivnet.GetEULA(tile.Product)
		if err != nil {
			return err
		}

		eulas[eula.Slug] = eula

		eula, err = cmd.pivnet.GetEULA(tile.Stemcell)
		if err != nil {
			return err
		}

		eulas[eula.Slug] = eula
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
	for _, tile := range cmd.pattern.Tiles {
		if err := cmd.pivnet.AcceptEULA(tile.Product); err != nil {
			return fmt.Errorf("accepting EULA for %s: %v", tile.Product.Slug, err)
		}

		if err := cmd.pivnet.AcceptEULA(tile.Stemcell); err != nil {
			return fmt.Errorf("accepting EULA for %s: %v", tile.Stemcell.Slug, err)
		}
	}

	cmd.logger.Printf("accepted EULAs for %d products", len(cmd.pattern.Tiles))

	return nil
}
