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
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"omg-cli/config"
	"omg-cli/ops_manager"
	"omg-cli/version"

	"github.com/alecthomas/kingpin"
)

type DirectorSSHCommand struct {
	logger *log.Logger
	envDir string
}

const DirectorSSHName = "director-ssh"

const jsonTemplate = `{
	"method": "ssh",
	"arguments": ["setup", {"user": "omg", "public-key": %q}]
}`

func (cmd *DirectorSSHCommand) register(app *kingpin.Application) {
	c := app.Command(DirectorSSHName, "Add the 'jumpbox' user and credentials to the BOSH director VM.").Action(cmd.run)
	registerEnvConfigFlag(c, &cmd.envDir)
}

func (cmd *DirectorSSHCommand) run(c *kingpin.ParseContext) error {
	cfg, err := config.TerraformFromEnvDirectory(cmd.envDir)
	if err != nil {
		return err
	}

	omSdk, err := ops_manager.NewSdk(fmt.Sprintf("https://%s", cfg.OpsManagerHostname), cfg.OpsManager, *cmd.logger)
	if err != nil {
		return err
	}

	cred, err := omSdk.GetDirectorCredentials("agent_credentials")
	if err != nil {
		return err
	}

	directorIP, err := omSdk.GetDirectorIP()
	if err != nil {
		return err
	}

	pubkey, err := ioutil.ReadFile(filepath.Join(cmd.envDir, "keys", "jumpbox_ssh.pub"))
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s:6868/agent", directorIP),
		bytes.NewBufferString(fmt.Sprintf(jsonTemplate, pubkey)))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", version.UserAgent())
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(cred.Identity, cred.Password)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cfg.OpsManager.SkipSSLVerification,
			},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Got %s from metron agent. Response body:\n%s\n", resp.Status, body)
	}
	return nil
}
