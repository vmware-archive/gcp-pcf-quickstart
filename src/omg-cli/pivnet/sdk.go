package pivnet

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"log"

	"omg-cli/tiles"

	"github.com/pivotal-cf/om/progress"
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
	req, err := http.NewRequest("GET", fmt.Sprintf("https://network.pivotal.io/%s", path), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", s.apiToken))

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
		return fmt.Errorf("authorizing pivnet credentials, recieved: %i", resp.StatusCode)
	}

	return nil
}

// DownloadTile retrieves a given productSlug, releaseId, and fileId from PivNet
// If a os.File is return it is guarenteed to match the fileSha256
// If an error is returned no os.File will be returned
//
// Caller is responsible for deleting the os.File
func (s *Sdk) DownloadTile(tile tiles.PivnetDefinition) (*os.File, error) {
	req, err := s.authorizedRequest("GET", fmt.Sprintf("/api/v2/products/%s/releases/%s/product_files/%s/download", tile.Name, tile.VersionId, tile.FileId), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	out, err := ioutil.TempFile("", "tile")
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

func (s *Sdk) AcceptEula(tile tiles.PivnetDefinition) error {
	req, err := s.authorizedRequest("POST", fmt.Sprintf("/api/v2/products/%s/releases/%s/eula_acceptance", tile.Name, tile.VersionId), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("accepting eula for %s, %s, recieved: %i", tile.Name, tile.VersionId, resp.StatusCode)
	}

	return nil
}
