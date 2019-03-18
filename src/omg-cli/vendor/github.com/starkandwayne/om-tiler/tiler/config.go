package tiler

import (
	"fmt"

	validator "gopkg.in/go-playground/validator.v9"
)

type Config struct {
	Target               string `yaml:"target" validate:"required"`
	Username             string `yaml:"username" validate:"required"`
	Password             string `yaml:"password" validate:"required"`
	DecryptionPassphrase string `yaml:"decryption_passphrase" validate:"required"`
	SkipSSLVerification  bool   `yaml:"skip_ssl_verification"`
	PivnetToken          string `yaml:"pivnet_token" validate:"required"`
}

func (c *Config) Validate() error {
	err := validator.New().Struct(c)
	if err != nil {
		return fmt.Errorf("tiler.Config has error(s):\n%+v\n", err)
	}
	return nil
}
