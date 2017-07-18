package app

import (
	"errors"
	"fmt"
	"log"
	"omg-cli/config"
	"omg-cli/omg/setup"
	"omg-cli/omg/tiles"
	"omg-cli/omg/tiles/ert"
	"omg-cli/ops_manager"
	"omg-cli/pivnet"
)

type App struct {
	setupService *setup.Service
}

func New(logger *log.Logger, mode Mode, terraformConfigPath string, pivnetAPIToken string, creds config.OpsManagerCredentials) (*App, error) {
	cfg, err := config.FromTerraform(terraformConfigPath)
	if err != nil {
		return nil, err
	}

	omSdk, err := ops_manager.NewSdk(fmt.Sprintf("https://%s", cfg.OpsManagerIp), creds, *logger)
	if err != nil {
		return nil, err
	}
	var pivnetSdk *pivnet.Sdk

	switch mode {
	case BakeImage:
		if pivnetAPIToken == "" {
			return nil, errors.New("pivnet-api-token required")
		}
		pivnetSdk, err = pivnet.NewSdk(pivnetAPIToken, logger)
		if err != nil {
			return nil, err
		}
	}

	selectedTiles := []tiles.TileInstaller{
		ert.Tile{},
	}

	return &App{setup.NewService(cfg, omSdk, pivnetSdk, logger, selectedTiles)}, nil
}

type step func() error

func runSteps(steps []step) error {
	for _, v := range steps {
		if err := v(); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) Run(mode Mode) error {
	switch mode {
	case BakeImage:
		return runSteps([]step{
			a.setupService.PoolTillOnline,
			a.setupService.SetupAuth,
			a.setupService.UploadTiles,
		})
	case ConfigureOpsManager:
		return runSteps([]step{
			a.setupService.PoolTillOnline,
			a.setupService.Unlock,
			//TODO(jrjohnson): RollCredentials
			a.setupService.SetupBosh,
			a.setupService.ConfigureTiles,
			//a.setupService.ApplyChanges,
		})
	default:
		return errors.New("unknown mode")
	}
}
