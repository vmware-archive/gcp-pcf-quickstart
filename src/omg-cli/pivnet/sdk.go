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
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"log"

	"omg-cli/config"

	"omg-cli/version"

	"time"

	"github.com/pivotal-cf/om/progress"
)

const (
	retryAttempts = 5 // How many times to retry downloading a tile from PivNet
	retryDelay    = 5 // How long wait in between download retries
)

type Sdk struct {
	apiToken string
	logger   *log.Logger
}

func NewSdk(apiToken string, logger *log.Logger) (*Sdk, error) {
	sdk := &Sdk{apiToken, logger}

	return sdk, sdk.checkCredentials()
}

func (s *Sdk) authorizedRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("https://network.pivotal.io/%s", path), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", s.apiToken))
	req.Header.Set("User-Agent", version.UserAgent())

	return req, nil
}

func (s *Sdk) checkCredentials() error {
	req, err := s.authorizedRequest("GET", "api/v2/authentication", nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("authorizing pivnet-api-token, recieved: %s", resp.Status)
	}

	return nil
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
	req, err := s.authorizedRequest("GET", fmt.Sprintf("/api/v2/products/%s/releases/%s/product_files/%s/download", tile.Name, tile.ReleaseId, tile.FileId), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("received non-200 %s response from server, body: %q", resp.Status, body)
	}

	var out *os.File
	if path == "" {
		out, err = ioutil.TempFile("", "tile")
	} else {
		out, err = os.Create(path)
	}
	if err != nil {
		return nil, err
	}

	// Delete the temp file if we're returning an error
	defer func() {
		if err != nil {
			os.Remove(out.Name())
		}
	}()

	// Stream the download to write the data to file in chunks
	// and calculate the sha256 the file is being written.
	//
	// This is necessary because we will read large (6+ GB) files
	// with this method.
	//
	// resp.body ==> BarReader (bar) ==> TeeReader ==> hasher (sha256)
	//                                             ==> out (temp file)
	s.logger.Printf("downloading tile: %s", tile.Name)
	hasher := sha256.New()
	bar := progress.NewBar()
	bar.SetTotal(resp.ContentLength)
	bar.Kickoff()
	defer bar.End()

	_, err = io.Copy(out, io.TeeReader(bar.NewBarReader(resp.Body), hasher))
	if err != nil {
		return nil, err
	}

	downloadedSha := fmt.Sprintf("%x", hasher.Sum(nil))
	if downloadedSha != tile.Sha256 {
		return nil, fmt.Errorf("sha256 of downloaded product does not match expected, got: %s, expected: %s", downloadedSha, tile.Sha256)
	}

	return out, nil
}

func (s *Sdk) AcceptEula(tile config.PivnetMetadata) error {
	req, err := s.authorizedRequest("POST", fmt.Sprintf("/api/v2/products/%s/releases/%s/eula_acceptance", tile.Name, tile.ReleaseId), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("accepting eula for %s, %s, recieved: %s", tile.Name, tile.ReleaseId, resp.Status)
	}

	return nil
}

type Eula struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	Slug    string `json:"slug"`
	Links   struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
}

func (s *Sdk) GetEula(eulaSlug string) (*Eula, error) {
	req, err := s.authorizedRequest("GET", "/api/v2/eulas", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetching eula: %s, recieved: %s", eulaSlug, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var eulaResponse struct {
		Eulas []Eula `json:"eulas"`
	}
	if err := json.Unmarshal(body, &eulaResponse); err != nil {
		return nil, fmt.Errorf("unmarshalling pivnet response: %v", err)
	}

	for _, eula := range eulaResponse.Eulas {
		if eula.Slug == eulaSlug {
			return &eula, nil
		}
	}

	return nil, fmt.Errorf("EULA not found: %s", eulaSlug)
}
