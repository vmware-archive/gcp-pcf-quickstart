package opsman

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/gosuri/uilive"
	"github.com/pivotal-cf/om/api"
	"github.com/pivotal-cf/om/commands"
	"github.com/pivotal-cf/om/extractor"
	"github.com/pivotal-cf/om/formcontent"
	"github.com/pivotal-cf/om/network"
	"github.com/pivotal-cf/om/progress"
	"github.com/starkandwayne/om-tiler/pattern"
)

type Config struct {
	Target               string
	Username             string
	Password             string
	DecryptionPassphrase string
	SkipSSLVerification  bool
}

type Client struct {
	api    api.Api
	log    *log.Logger
	config Config
}

const (
	connectTimeout     = time.Duration(5) * time.Second
	requestTimeout     = time.Duration(1800) * time.Second
	pollingIntervalSec = "10"
)

func NewClient(c Config, logger *log.Logger) (*Client, error) {
	oauthClient, err := network.NewOAuthClient(
		c.Target, c.Username, c.Password, "", "",
		c.SkipSSLVerification, true,
		requestTimeout, connectTimeout,
	)
	if err != nil {
		return &Client{}, err
	}

	unauthenticatedClient := network.NewUnauthenticatedClient(
		c.Target, c.SkipSSLVerification,
		requestTimeout, connectTimeout,
	)

	logger.SetPrefix(fmt.Sprintf("%s[OM] ", logger.Prefix()))

	live := uilive.New()
	live.Out = os.Stderr

	client := Client{
		api: api.New(api.ApiInput{
			Client:         oauthClient,
			UnauthedClient: unauthenticatedClient,
			ProgressClient: network.NewProgressClient(
				oauthClient, progress.NewBar(), live),
			UnauthedProgressClient: network.NewProgressClient(
				unauthenticatedClient, progress.NewBar(), live),
			Logger: logger,
		}),
		log:    logger,
		config: c,
	}

	return &client, nil
}

func (c *Client) ConfigureAuthentication() error {
	args := []string{
		fmt.Sprintf("--username=%s", c.config.Username),
		fmt.Sprintf("--password=%s", c.config.Password),
		fmt.Sprintf("--decryption-passphrase=%s", c.config.DecryptionPassphrase),
	}
	cmd := commands.NewConfigureAuthentication(c.api, c.log)
	return cmd.Execute(args)
}

func (c *Client) FilesUploaded(t pattern.Tile) (bool, error) {
	return false, nil
}

func (c *Client) UploadProduct(p *os.File) error {
	args := []string{
		fmt.Sprintf("--product=%s", p.Name()),
		fmt.Sprintf("--polling-interval=%s", pollingIntervalSec),
	}
	form := formcontent.NewForm()
	metadataExtractor := extractor.MetadataExtractor{}
	cmd := commands.NewUploadProduct(form, metadataExtractor, c.api, c.log)
	return cmd.Execute(args)
}

func (c *Client) UploadStemcell(s *os.File) error {
	args := []string{
		fmt.Sprintf("--stemcell=%s", s.Name()),
		"--floating",
	}
	form := formcontent.NewForm()
	cmd := commands.NewUploadStemcell(form, c.api, c.log)
	return cmd.Execute(args)
}

func (c *Client) StageProduct(t pattern.Tile) error {
	args := []string{
		fmt.Sprintf("--product-name=%s", t.Name),
		fmt.Sprintf("--product-version=%s", t.Version),
	}
	cmd := commands.NewStageProduct(c.api, c.log)
	return cmd.Execute(args)
}

func (c *Client) ConfigureProduct(config []byte) error {
	configFile, err := tmpConfigFile(config)
	if err != nil {
		return err
	}
	defer os.Remove(configFile)

	args := []string{
		fmt.Sprintf("--config=%s", configFile),
	}
	cmd := commands.NewConfigureProduct(
		os.Environ, c.api, c.config.Target, c.log)
	return cmd.Execute(args)
}

func (c *Client) ConfigureDirector(config []byte) error {
	configFile, err := tmpConfigFile(config)
	if err != nil {
		return err
	}
	defer os.Remove(configFile)

	args := []string{
		fmt.Sprintf("--config=%s", configFile),
	}
	cmd := commands.NewConfigureDirector(os.Environ, c.api, c.log)
	return cmd.Execute(args)
}

func (c *Client) ApplyChanges() error {
	return nil
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
