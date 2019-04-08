package opsman

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	api                   api.Api
	log                   *log.Logger
	config                Config
	unauthenticatedClient network.UnauthenticatedClient
}

const (
	connectTimeout     = time.Duration(5) * time.Second
	requestTimeout     = time.Duration(1800) * time.Second
	pollingInterval    = time.Duration(10) * time.Second
	applySleepDuration = time.Duration(10) * time.Second
	onlineTimeout      = time.Duration(240 * time.Second)
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
		log:                   logger,
		config:                c,
		unauthenticatedClient: unauthenticatedClient,
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
	products, err := c.uploadedProducts()
	if err != nil {
		return false, err
	}

	pok := contains(products, fmt.Sprintf("%s/%s", t.Name, t.Version))
	sok := contains(products, fmt.Sprintf("stemcell/%s", t.Stemcell.Version))

	return (pok && sok), nil
}

func (c *Client) UploadProduct(p *os.File) error {
	args := []string{
		fmt.Sprintf("--product=%s", p.Name()),
		fmt.Sprintf("--polling-interval=%d", int(pollingInterval.Seconds())),
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
	args := []string{"--skip-unchanged-products"}
	logWriter := commands.NewLogWriter(os.Stdout)
	cmd := commands.NewApplyChanges(c.api, c.api, logWriter, c.log, applySleepDuration)
	return cmd.Execute(args)
}

func (c *Client) DeleteInstallation() error {
	logWriter := commands.NewLogWriter(os.Stdout)
	cmd := commands.NewDeleteInstallation(c.api, logWriter, c.log, pollingInterval)
	return cmd.Execute(nil)
}

func (c *Client) PollTillOnline() error {
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
			c.log.Print("waiting for Ops Manager to start")
			timer = time.After(pollingInterval)
		}
	}
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

type StemcellAssignments struct {
	Products []Product `json:"products"`
}

type Product struct {
	Name             string   `json:"type"`
	Version          string   `json:"product_version"`
	StemcellVersions []string `json:"available_stemcell_versions"`
}

func (c *Client) uploadedProducts() ([]string, error) {
	args := []string{"--silent", "--path=/api/v0/stemcell_assignments"}
	out := bytes.NewBuffer([]byte{})
	cmd := commands.NewCurl(c.api, log.New(out, "", 0), c.log)
	err := cmd.Execute(args)
	if err != nil {
		return []string{}, fmt.Errorf("retrieving stemcell assignments: %s", err)
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
