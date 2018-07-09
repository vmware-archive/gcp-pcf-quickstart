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

package pivnet

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"omg-cli/config"
	"omg-cli/version"

	gopivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/logshim"
)

const (
	retryAttempts = 5 // How many times to retry downloading a tile from PivNet
	retryDelay    = 5 // How long wait in between download retries
)

type Sdk struct {
	logger   *log.Logger
	client   gopivnet.Client
	apiToken string
}

func NewSdk(apiToken string, logger *log.Logger) (*Sdk, error) {
	sdk := &Sdk{logger: logger, apiToken: apiToken}

	cfg := gopivnet.ClientConfig{
		Host:      gopivnet.DefaultHost,
		Token:     apiToken,
		UserAgent: version.UserAgent(),
	}
	sdk.client = gopivnet.NewClient(cfg, logshim.NewLogShim(logger, logger, false))

	return sdk, sdk.checkCredentials()
}

func (s *Sdk) checkCredentials() error {
	ok, err := s.client.Auth.Check()

	if ok {
		return nil
	} else {
		return fmt.Errorf("authorizing pivnet credentials: %v", err)
	}
}

// DownloadTile retrieves a given productSlug, releaseId, and fileId from PivNet
// If a os.File is return it is guarenteed to match the fileSha256
// If an error is returned no os.File will be returned
//
// Caller is responsible for deleting the os.File
func (s *Sdk) DownloadTile(tile config.PivnetMetadata) (file *os.File, err error) {
	return s.DownloadTileToPath(tile, "")
}

// DownloadTileToPath is a version of DownloadTile that accepts a path specifying
// download location of the tile.
func (s *Sdk) DownloadTileToPath(tile config.PivnetMetadata, path string) (file *os.File, err error) {
	for i := 0; i < retryAttempts; i++ {
		file, err = s.downloadTile(tile, path)

		// Success or recoverable error
		if err == nil || err != io.ErrUnexpectedEOF {
			return
		}

		s.logger.Printf("download tile failed, retrying in %d seconds", retryDelay)
		time.Sleep(time.Duration(retryDelay) * time.Second)
	}

	return nil, fmt.Errorf("download tile failed after %d attempts", retryAttempts)
}

func (s *Sdk) downloadTile(tile config.PivnetMetadata, path string) (*os.File, error) {
	var err error
	var out *os.File
	if path == "" {
		out, err = ioutil.TempFile("", "tile")
	} else {
		out, err = os.Create(path)
	}
	if err != nil {
		return nil, err
	}

	// Delete the file if we're returning an error
	defer func() {
		if err != nil {
			os.Remove(out.Name())
		}
	}()

	return out, s.client.ProductFiles.DownloadForRelease(out, tile.Name, tile.ReleaseId, tile.FileId, os.Stdout)
}

func (s *Sdk) AcceptEula(tile config.PivnetMetadata) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("https://network.pivotal.io/api/v2/products/%s/releases/%d/eula_acceptance", tile.Name, tile.ReleaseId), nil)
	if err != nil {
		return err
	}
	response, err := s.client.Auth.FetchUAAToken(s.apiToken)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", response.Token))
	req.Header.Set("User-Agent", version.UserAgent())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("accepting eula for %s, %v, recieved: %s", tile.Name, tile.ReleaseId, resp.Status)
	}

	return nil
}

type Eula struct {
	Name    string
	Content string
	Slug    string
}

func (s *Sdk) GetEula(eulaSlug string) (*Eula, error) {
	eula, err := s.client.EULA.Get(eulaSlug)
	if err != nil {
		return nil, err
	}

	return &Eula{Name: eula.Name, Content: eula.Content, Slug: eula.Slug}, nil
}
