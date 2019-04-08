package pivnet

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	gopivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/logshim"
	"github.com/pivotal-cf/pivnet-cli/filter"
	"github.com/starkandwayne/om-tiler/pattern"
)

const (
	retryAttempts = 5 // How many times to retry downloading a tile from PivNet
	retryDelay    = 5 // How long wait in between download retries
)

type Config struct {
	Host       string
	Token      string
	UserAgent  string
	AcceptEULA bool
}

type Client struct {
	logger     *log.Logger
	client     gopivnet.Client
	acceptEULA bool
	userAgent  string
	filter     *filter.Filter
}

type EULA struct {
	Name    string
	Content string
	Slug    string
}

func NewClient(c Config, logger *log.Logger) *Client {
	host := c.Host
	if c.Host == "" {
		host = gopivnet.DefaultHost
	}
	log := logshim.NewLogShim(logger, logger, false)
	client := gopivnet.NewClient(gopivnet.ClientConfig{
		Host:      host,
		Token:     c.Token,
		UserAgent: c.UserAgent,
	}, log)
	filter := filter.NewFilter(log)
	return &Client{client: client, logger: logger,
		acceptEULA: c.AcceptEULA, userAgent: c.UserAgent, filter: filter}
}

func (c *Client) DownloadFile(f pattern.PivnetFile, dir string) (file *os.File, err error) {
	if c.acceptEULA {
		if err = c.AcceptEULA(f); err != nil {
			return
		}
	}
	for i := 0; i < retryAttempts; i++ {
		file, err = c.downloadFile(f, dir)

		// Success or recoverable error
		if err == nil || err != io.ErrUnexpectedEOF {
			return
		}

		c.logger.Printf("download tile failed, retrying in %d seconds", retryDelay)
		time.Sleep(time.Duration(retryDelay) * time.Second)
	}

	return nil, fmt.Errorf("download tile failed after %d attempts", retryAttempts)
}

func (c *Client) GetEULA(f pattern.PivnetFile) (*EULA, error) {
	release, err := c.lookupRelease(f)
	if err != nil {
		return nil, err
	}

	eula, err := c.client.EULA.Get(release.EULA.Slug)
	if err != nil {
		return nil, err
	}

	return &EULA{Name: eula.Name, Content: eula.Content, Slug: eula.Slug}, nil
}

func (c *Client) AcceptEULA(f pattern.PivnetFile) error {
	release, err := c.lookupRelease(f)
	if err != nil {
		return err
	}

	if c.userAgent != "" {
		return c.forceAcceptEULA(f.Slug, release.ID)
	}
	return c.client.EULA.Accept(f.Slug, release.ID)
}

func (c *Client) forceAcceptEULA(productSlug string, releaseID int) error {
	url := fmt.Sprintf(
		"/products/%s/releases/%d/eula_acceptance",
		productSlug,
		releaseID,
	)

	resp, err := c.client.MakeRequest(
		"POST",
		url,
		http.StatusOK,
		strings.NewReader(`{}`),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) downloadFile(f pattern.PivnetFile, dir string) (file *os.File, err error) {
	if dir == "" {
		dir, err = ioutil.TempDir("", f.Slug)
		if err != nil {
			return nil, err
		}
	}

	productFile, release, err := c.lookupProductFile(f)
	if err != nil {
		return nil, err
	}

	baseName := filepath.Base(productFile.AWSObjectKey)
	file, err = os.Create(filepath.Join(dir, baseName))
	if err != nil {
		return nil, err
	}

	// Delete the file if we're returning an error
	defer func() {
		if err != nil {
			os.RemoveAll(dir)
		}
	}()

	return file, c.client.ProductFiles.DownloadForRelease(file, f.Slug, release.ID, productFile.ID, os.Stdout)
}

func (c *Client) lookupRelease(f pattern.PivnetFile) (gopivnet.Release, error) {
	releases, err := c.client.Releases.List(f.Slug)
	if err != nil {
		return gopivnet.Release{}, err
	}

	for _, r := range releases {
		if r.Version == f.Version {
			return r, nil
		}
	}
	return gopivnet.Release{}, fmt.Errorf(
		"release not found for %s with version: '%s'", f.Slug, f.Version,
	)
}

func (c *Client) lookupProductFile(f pattern.PivnetFile) (gopivnet.ProductFile, gopivnet.Release, error) {
	release, err := c.lookupRelease(f)
	productFiles, err := c.client.ProductFiles.ListForRelease(f.Slug, release.ID)
	if err != nil {
		return gopivnet.ProductFile{}, gopivnet.Release{}, err
	}

	productFiles, err = c.filter.ProductFileKeysByGlobs(productFiles, []string{f.Glob})
	if err != nil {
		return gopivnet.ProductFile{}, gopivnet.Release{},
			fmt.Errorf("could not glob product files: %s", err)
	}

	if err := c.checkForSingleProductFile(f.Glob, productFiles); err != nil {
		return gopivnet.ProductFile{}, gopivnet.Release{}, err
	}

	return productFiles[0], release, nil

}

func (c *Client) checkForSingleProductFile(glob string, productFiles []gopivnet.ProductFile) error {
	if len(productFiles) > 1 {
		var productFileNames []string
		for _, productFile := range productFiles {
			productFileNames = append(productFileNames, path.Base(productFile.AWSObjectKey))
		}
		return fmt.Errorf("the glob '%s' matches multiple files. Write your glob to match exactly one of the following:\n  %s", glob, strings.Join(productFileNames, "\n  "))
	} else if len(productFiles) == 0 {
		return fmt.Errorf("the glob '%s' matches no file", glob)
	}

	return nil
}
