package opsman

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pivotal-cf/om/api"
	"github.com/pivotal-cf/om/commands"
	"github.com/pivotal-cf/om/extractor"
	"github.com/pivotal-cf/om/formcontent"
	"github.com/pivotal-cf/om/network"
	"github.com/starkandwayne/om-tiler/pattern"
	"github.com/starkandwayne/om-tiler/steps"
)

// Config information needed to create a Client
type Config struct {
	Target               string
	Username             string
	Password             string
	DecryptionPassphrase string
	SkipSSLVerification  bool
}

// Client responsible for interacting with the OpsManager API
type Client struct {
	config                Config
	logger                func(context.Context) *log.Logger
	unauthenticatedClient network.UnauthenticatedClient
	oauthClient           network.OAuthClient
}

const (
	connectTimeout     = time.Duration(5) * time.Second
	requestTimeout     = time.Duration(1800) * time.Second
	pollingInterval    = time.Duration(10) * time.Second
	applySleepDuration = time.Duration(10) * time.Second
	onlineTimeout      = time.Duration(240 * time.Second)
)

// NewClient creates a new Client when given Config and a logger
// will attempt to get current Step from context.Context to prefix logger
func NewClient(c Config, logger *log.Logger) (*Client, error) {
	oauthClient, err := network.NewOAuthClient(
		c.Target, c.Username, c.Password, "", "",
		c.SkipSSLVerification, true,
		requestTimeout, connectTimeout,
	)
	if err != nil {
		return &Client{}, fmt.Errorf(
			"could not create OpsManager OAuth Client: %v", err)
	}

	log := func(ctx context.Context) *log.Logger {
		return steps.ContextLogger(ctx, logger, "[OM]")
	}

	return &Client{
		config:      c,
		logger:      log,
		oauthClient: oauthClient,
		unauthenticatedClient: network.NewUnauthenticatedClient(
			c.Target, c.SkipSSLVerification,
			requestTimeout, connectTimeout,
		),
	}, nil
}

// ConfigureAuthentication configures OpsManager authentication
func (c *Client) ConfigureAuthentication(ctx context.Context) error {
	args := []string{
		fmt.Sprintf("--username=%s", c.config.Username),
		fmt.Sprintf("--password=%s", c.config.Password),
		fmt.Sprintf("--decryption-passphrase=%s", c.config.DecryptionPassphrase),
	}
	cmd := commands.NewConfigureAuthentication(c.api(ctx), c.logger(ctx))
	return cmd.Execute(args)
}

// FilesUploaded checks if all PivnetFiles for a given Tile have been uploaded
func (c *Client) FilesUploaded(ctx context.Context, t pattern.Tile) (bool, error) {
	products, err := c.uploadedProducts(ctx)
	if err != nil {
		return false, fmt.Errorf("could not fetch already uploaded products: %v", err)
	}

	pok := contains(products, fmt.Sprintf("%s/%s", t.Name, t.Version))
	sok := contains(products, fmt.Sprintf("stemcell/%s", t.Stemcell.Version))

	return (pok && sok), nil
}

// UploadProduct uploads a given tile file to OpsManager
func (c *Client) UploadProduct(ctx context.Context, p *os.File) error {
	args := []string{
		fmt.Sprintf("--product=%s", p.Name()),
		fmt.Sprintf("--polling-interval=%d", int(pollingInterval.Seconds())),
	}
	form := formcontent.NewForm()
	metadataExtractor := extractor.MetadataExtractor{}
	cmd := commands.NewUploadProduct(form, metadataExtractor, c.api(ctx), c.logger(ctx))
	return cmd.Execute(args)
}

// UploadStemcell uploads a given stemcell file to OpsManager
func (c *Client) UploadStemcell(ctx context.Context, s *os.File) error {
	args := []string{
		fmt.Sprintf("--stemcell=%s", s.Name()),
		"--floating=true",
	}
	form := formcontent.NewForm()
	cmd := commands.NewUploadStemcell(form, c.api(ctx), c.logger(ctx))
	return cmd.Execute(args)
}

// StageProduct stages the product with Tile.Name and Tile.Version
func (c *Client) StageProduct(ctx context.Context, t pattern.Tile) error {
	args := []string{
		fmt.Sprintf("--product-name=%s", t.Name),
		fmt.Sprintf("--product-version=%s", t.Version),
	}
	cmd := commands.NewStageProduct(c.api(ctx), c.logger(ctx))
	return cmd.Execute(args)
}

// ConfigureProduct configures a product when given an om cli compatible manifest
func (c *Client) ConfigureProduct(ctx context.Context, config []byte) error {
	configFile, err := tmpConfigFile(config)
	if err != nil {
		return fmt.Errorf("could not create temporary product config file: %v", err)
	}
	defer os.Remove(configFile)

	args := []string{
		fmt.Sprintf("--config=%s", configFile),
	}
	cmd := commands.NewConfigureProduct(
		os.Environ, c.api(ctx), c.config.Target, c.logger(ctx))
	return cmd.Execute(args)
}

// ConfigureDirector configures the director when given an om cli compatible manifest
func (c *Client) ConfigureDirector(ctx context.Context, config []byte) error {
	configFile, err := tmpConfigFile(config)
	if err != nil {
		return fmt.Errorf("could not create temporary director config file: %v", err)
	}
	defer os.Remove(configFile)

	args := []string{
		fmt.Sprintf("--config=%s", configFile),
	}
	cmd := commands.NewConfigureDirector(os.Environ, c.api(ctx), c.logger(ctx))
	return cmd.Execute(args)
}

// ApplyChanges applies all pending changes (will skip unchanged products)
func (c *Client) ApplyChanges(ctx context.Context) error {
	args := []string{}
	logWriter := commands.NewLogWriter(os.Stdout)
	api := c.api(ctx)
	cmd := commands.NewApplyChanges(api, api, logWriter, c.logger(ctx), applySleepDuration)
	return cmd.Execute(args)
}

// DeleteInstallation deletes all tiles and their resources
func (c *Client) DeleteInstallation(ctx context.Context) error {
	args := []string{"--force"}
	logWriter := commands.NewLogWriter(os.Stdout)
	cmd := commands.NewDeleteInstallation(c.api(ctx), logWriter, c.logger(ctx), os.Stdin, pollingInterval)
	return cmd.Execute(args)
}

// PollTillOnline wait for the OpsManager API to become available
func (c *Client) PollTillOnline(ctx context.Context) error {
	timer := time.After(time.Duration(0 * time.Second))
	timeout := time.After(onlineTimeout)
	for {
		select {
		case <-timeout:
			return errors.New("timeout waiting for Ops Manager to start")
		case <-timer:
			if c.online() {
				return nil
			}
			c.logger(ctx).Print("waiting for Ops Manager to start")
			timer = time.After(pollingInterval)
		}
	}
}

func (c *Client) api(ctx context.Context) api.Api {
	live := liveDiscarder{}

	return api.New(api.ApiInput{
		Client:         c.oauthClient,
		UnauthedClient: c.unauthenticatedClient,
		ProgressClient: network.NewProgressClient(
			c.oauthClient, progressBarDiscarder{}, live),
		UnauthedProgressClient: network.NewProgressClient(
			c.unauthenticatedClient, progressBarDiscarder{}, live),
		Logger: c.logger(ctx),
	})
}

func (c *Client) online() bool {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	resp, err := c.unauthenticatedClient.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode < 500
}

func tmpConfigFile(config []byte) (string, error) {
	configFile, err := ioutil.TempFile("", "config")
	if err != nil {
		return "", err
	}

	if _, err = configFile.Write(config); err != nil {
		return "", err
	}

	if err = configFile.Close(); err != nil {
		return "", err
	}

	return configFile.Name(), nil
}

func (c *Client) uploadedProducts(ctx context.Context) ([]string, error) {
	args := []string{"--silent", "--path=/api/v0/stemcell_assignments"}
	out := bytes.NewBuffer([]byte{})
	cmd := commands.NewCurl(c.api(ctx), log.New(out, "", 0), c.logger(ctx))
	err := cmd.Execute(args)
	if err != nil {
		return []string{}, fmt.Errorf("retrieving stemcell assignments: %s", err)
	}

	type Product struct {
		Name             string   `json:"identifier"`
		Version          string   `json:"staged_product_version"`
		StemcellVersions []string `json:"available_stemcell_versions"`
	}

	type StemcellAssignments struct {
		Products []Product `json:"products"`
	}

	var assignments StemcellAssignments
	err = json.Unmarshal(out.Bytes(), &assignments)
	if err != nil {
		return []string{}, fmt.Errorf("decoding stemcell assignments: %s", err)
	}

	products := []string{}
	for _, p := range assignments.Products {
		products = append(products, fmt.Sprintf("%s/%s", p.Name, p.Version))
		for _, s := range p.StemcellVersions {
			products = append(products, fmt.Sprintf("stemcell/%s", s))
		}
	}
	return products, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
