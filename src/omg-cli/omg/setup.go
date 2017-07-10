package omg

import (
	"omg-cli/config"
	"omg-cli/ops_manager"

	"github.com/pivotal-cf/om/commands"
)

type SetupService struct {
	cfg *config.Config
	sdk *ops_manager.Sdk
}

func NewSetupService(cfg *config.Config, sdk *ops_manager.Sdk) *SetupService {
	return &SetupService{cfg: cfg, sdk: sdk}
}

func (s *SetupService) SetupAuth(decryptionPhrase string) error {
	return s.sdk.SetupAuth(decryptionPhrase)
}

func (s *SetupService) SetupBosh() error {
	iassCfg := commands.GCPIaaSConfiguration{
		Project:              s.cfg.ProjectName,
		DefaultDeploymentTag: "omg-opsman",
		AuthJSON:             "",
	}

	if err := s.sdk.SetupBosh(iassCfg); err != nil {
		return err
	}

	return nil
}
