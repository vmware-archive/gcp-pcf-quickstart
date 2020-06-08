package pivnet

import (
	"context"
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
	"github.com/pivotal-cf/go-pivnet/download"
	"github.com/pivotal-cf/go-pivnet/logshim"
	"github.com/pivotal-cf/pivnet-cli/filter"
	"github.com/starkandwayne/om-tiler/pattern"
	"github.com/starkandwayne/om-tiler/steps"
)

const (
	retryAttempts = 5 // How many times to retry downloading a tile from PivNet
	retryDelay    = 5 // How long wait in between download retries
)

// Config information needed to create a Client
type Config struct {
	Host       string
	Token      string
	UserAgent  string
	AcceptEULA bool
}

// Client responsible for interacting with Pivotal Network API
type Client struct {
	acceptEULA bool
	userAgent  string
	logger     func(context.Context) *log.Logger
	client     func(context.Context) gopivnet.Client
	filter     func(context.Context) *filter.Filter
}

// EULA contains text of a Pivnet End User License Agreement
type EULA struct {
	Name    string
	Content string
	Slug    string
}

// NewClient creates a new Client when given Config and a logger
// will attempt to get current Step from context.Context to prefix logger
func NewClient(c Config, logger *log.Logger) *Client {
	host := c.Host
	if c.Host == "" {
		host = gopivnet.DefaultHost
	}

	log := func(ctx context.Context) *log.Logger {
		return steps.ContextLogger(ctx, logger, "[Pivnet]")
	}
	client := func(ctx context.Context) gopivnet.Client {
		return gopivnet.NewClient(
			gopivnet.NewAccessTokenOrLegacyToken(c.Token, host, c.UserAgent),
			gopivnet.ClientConfig{Host: host, UserAgent: c.UserAgent},
			logshim.NewLogShim(log(ctx), log(ctx), false))
	}
	filter := func(ctx context.Context) *filter.Filter {
		return filter.NewFilter(logshim.NewLogShim(log(ctx), log(ctx), false))
	}
	return &Client{client: client, logger: log,
		acceptEULA: c.AcceptEULA, userAgent: c.UserAgent, filter: filter}
}

// DownloadFile downloads a PivnetFile into a destination directory
// will invoke AcceptEULA unless downloading form a URL
func (c *Client) DownloadFile(ctx context.Context, f pattern.PivnetFile, dir string) (file *os.File, err error) {
	if c.acceptEULA && f.URL == "" {
		if err = c.AcceptEULA(ctx, f); err != nil {
			return nil, fmt.Errorf("could not accept EULA: %v", err)
		}
	}
	for i := 0; i < retryAttempts; i++ {
		file, err = c.downloadFile(ctx, f, dir)

		// Success or recoverable error
		if err == nil || err != io.ErrUnexpectedEOF {
			return
		}

		c.logger(ctx).Printf("download tile failed, retrying in %d seconds", retryDelay)
		time.Sleep(time.Duration(retryDelay) * time.Second)
	}

	return nil, fmt.Errorf("download tile failed after %d attempts: %v", retryAttempts, err)
}

// GetEULA retrieves EULA for a given PivnetFile
func (c *Client) GetEULA(ctx context.Context, f pattern.PivnetFile) (*EULA, error) {
	release, err := c.lookupRelease(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("could not lookup release: %v", err)
	}

	eula, err := c.client(ctx).EULA.Get(release.EULA.Slug)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve EULA for %s: %v",
			release.EULA.Slug, err)
	}

	return &EULA{Name: eula.Name, Content: eula.Content, Slug: eula.Slug}, nil
}

// AcceptEULA accepts eula for a given PivnetFile
// will use different endpoint when client is configured with UserAgent
func (c *Client) AcceptEULA(ctx context.Context, f pattern.PivnetFile) error {
	release, err := c.lookupRelease(ctx, f)
	if err != nil {
		return fmt.Errorf("could not lookup release: %v", err)
	}

	if c.userAgent != "" {
		return c.forceAcceptEULA(ctx, f.Slug, release.ID)
	}
	return c.client(ctx).EULA.Accept(f.Slug, release.ID)
}

func (c *Client) forceAcceptEULA(ctx context.Context, productSlug string, releaseID int) error {
	url := fmt.Sprintf(
		"/products/%s/releases/%d/eula_acceptance",
		productSlug,
		releaseID,
	)

	resp, err := c.client(ctx).MakeRequest(
		"POST",
		url,
		http.StatusOK,
		strings.NewReader(`{}`),
	)
	if err != nil {
		return fmt.Errorf("POST to %s failed: %v", url, err)
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) downloadFile(ctx context.Context, f pattern.PivnetFile, dir string) (file *os.File, err error) {
	if dir == "" {
		dir, err = ioutil.TempDir("", f.Slug)
		if err != nil {
			return nil, fmt.Errorf("could not create tmp download dir: %v", err)
		}
	}

	// Delete the file if we're returning an error
	defer func() {
		if err != nil {
			os.RemoveAll(dir)
		}
	}()

	if f.URL != "" {
		file, err = downloadDirectFile(ctx, f.URL, dir)
		if err != nil {
			return nil, fmt.Errorf("could not download file from url %s: %v", f.URL, err)
		}
		return
	}
	return c.downloadPivnetFile(ctx, f, dir)
}

func (c *Client) downloadPivnetFile(ctx context.Context, f pattern.PivnetFile, dir string) (file *os.File, err error) {
	productFile, release, err := c.lookupProductFile(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("could not lookup product file: %v", err)
	}

	baseName := filepath.Base(productFile.AWSObjectKey)
	file, err = os.Create(filepath.Join(dir, baseName))
	if err != nil {
		return nil, fmt.Errorf("could not create download file %s: %v", file.Name(), err)
	}

	fi, err := download.NewFileInfo(file)
	if err != nil {
		return nil, fmt.Errorf("could not get file info %s: %v", file.Name(), err)
	}
	return file, c.client(ctx).ProductFiles.DownloadForRelease(fi, f.Slug, release.ID, productFile.ID, ioutil.Discard)
}

func downloadDirectFile(ctx context.Context, url string, dir string) (file *os.File, err error) {
	baseName := filepath.Base(url)
	file, err = os.Create(filepath.Join(dir, baseName))
	if err != nil {
		return nil, fmt.Errorf("could not create download file %s: %v", file.Name(), err)
	}
	defer file.Close()
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET on %s failed: %v", url, err)
	}

	defer resp.Body.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not copy response body: %v", err)
	}
	return
}

func (c *Client) lookupRelease(ctx context.Context, f pattern.PivnetFile) (gopivnet.Release, error) {
	releases, err := c.client(ctx).Releases.List(f.Slug)
	if err != nil {
		return gopivnet.Release{}, fmt.Errorf("failed to list releases: %v", err)
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

func (c *Client) lookupProductFile(ctx context.Context, f pattern.PivnetFile) (gopivnet.ProductFile, gopivnet.Release, error) {
	release, err := c.lookupRelease(ctx, f)
	if err != nil {
		return gopivnet.ProductFile{}, gopivnet.Release{}, fmt.Errorf(
			"could not find release for product file: %v", err)
	}

	productFiles, err := c.client(ctx).ProductFiles.ListForRelease(f.Slug, release.ID)
	if err != nil {
		return gopivnet.ProductFile{}, gopivnet.Release{}, fmt.Errorf(
			"could not list product files for release: %v", err)
	}

	productFiles, err = c.filter(ctx).ProductFileKeysByGlobs(productFiles, []string{f.Glob})
	if err != nil {
		return gopivnet.ProductFile{}, gopivnet.Release{},
			fmt.Errorf("could not glob product files: %s", err)
	}

	if err := c.checkForSingleProductFile(f.Glob, productFiles); err != nil {
		return gopivnet.ProductFile{}, gopivnet.Release{}, fmt.Errorf(
			"could not ensure single product file: %v", err)
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
