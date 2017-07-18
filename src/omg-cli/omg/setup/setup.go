package setup

import (
	"omg-cli/config"
	"omg-cli/omg/bosh_director"
	"omg-cli/ops_manager"
	"omg-cli/pivnet"

	"errors"
	"log"
	"omg-cli/omg/tiles"
	"os"
	"time"
)

type Service struct {
	cfg    *config.Config
	om     *ops_manager.Sdk
	pivnet *pivnet.Sdk
	logger *log.Logger
	tiles  []tiles.TileInstaller
}

func NewService(cfg *config.Config, omSdk *ops_manager.Sdk, pivnetSdk *pivnet.Sdk, logger *log.Logger, tiles []tiles.TileInstaller) *Service {
	return &Service{cfg, omSdk, pivnetSdk, logger, tiles}
}

func (s *Service) SetupAuth() error {
	return s.om.SetupAuth()
}

func (s *Service) Unlock() error {
	err := s.om.Unlock()
	if err != nil {
		return err
	}

	timer := time.After(time.Duration(0 * time.Second))
	timeout := time.After(time.Duration(120 * time.Second))
	for {
		select {
		case <-timeout:
			return errors.New("Timeout waiting for Ops Manager to unlock")
		case <-timer:
			if s.om.ReadyForAuth() {
				return nil
			}
			s.logger.Print("waiting for Ops Manager to unlock")
			timer = time.After(time.Duration(5 * time.Second))
		}
	}
}

func (s *Service) SetupBosh() error {
	gcp := bosh_director.GCP(s.cfg)
	director := bosh_director.Director()
	azs := bosh_director.AvalibilityZones(s.cfg)
	networks, networkAssignment := bosh_director.Network(s.cfg)
	resources := bosh_director.Resources()

	if err := s.om.SetupBosh(gcp, director, azs, networks, networkAssignment, resources); err != nil {
		return err
	}

	return nil
}

func (s *Service) ApplyChanges() error {
	return s.om.ApplyChanges()
}

func (s *Service) productInstalled(product config.OpsManagerMetadata) (bool, error) {
	products, err := s.om.AvaliableProducts()
	if err != nil {
		return false, err
	}

	for _, p := range products {
		if p.Name == product.Name && p.Version == product.Version {
			return true, nil
		}
	}
	return false, nil
}

func (s *Service) ensureProductReady(tile config.Tile) error {
	if i, err := s.productInstalled(tile.Product); i == true || err != nil {
		return err
	}

	file, err := s.pivnet.DownloadTile(tile.Pivnet)
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	if err = s.om.UploadProduct(file.Name()); err != nil {
		return err
	}

	return s.om.StageProduct(tile.Product)
}

func (s *Service) PoolTillOnline() error {
	timer := time.After(time.Duration(0 * time.Second))
	timeout := time.After(time.Duration(120 * time.Second))
	for {
		select {
		case <-timeout:
			return errors.New("Timeout waiting for Ops Manager to start")
		case <-timer:
			if s.om.Online() {
				return nil
			}
			s.logger.Print("waiting for Ops Manager to start")
			timer = time.After(time.Duration(5 * time.Second))
		}
	}
}

func (s *Service) ConfigureTiles() error {
	for _, t := range s.tiles {
		if err := t.Configure(s.cfg, s.om); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) UploadTiles() error {
	for _, t := range s.tiles {
		if err := s.ensureProductReady(t.Definition()); err != nil {
			return err
		}
	}

	return nil
}
