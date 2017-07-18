package omg

import (
	"omg-cli/config"
	"omg-cli/omg/bosh_director"
	"omg-cli/omg/ert"
	"omg-cli/omg/service_broker"
	"omg-cli/omg/stackdriver_nozzle"
	"omg-cli/ops_manager"
	"omg-cli/pivnet"
	"omg-cli/tiles"

	"errors"
	"log"
	"os"
	"time"
)

type SetupService struct {
	cfg    *config.Config
	om     *ops_manager.Sdk
	pivnet *pivnet.Sdk
	logger *log.Logger
}

func NewSetupService(cfg *config.Config, omSdk *ops_manager.Sdk, pivnetSdk *pivnet.Sdk, logger *log.Logger) *SetupService {
	return &SetupService{cfg, omSdk, pivnetSdk, logger}
}

func (s *SetupService) SetupAuth(decryptionPhrase string) error {
	return s.om.SetupAuth(decryptionPhrase)
}

func (s *SetupService) Unlock(decryptionPhrase string) error {
	err := s.om.Unlock(decryptionPhrase)
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

func (s *SetupService) SetupBosh() error {
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

func (s *SetupService) ApplyChanges() error {
	return s.om.ApplyChanges()
}

func (s *SetupService) productInstalled(product tiles.ProductDefinition) (bool, error) {
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

func (s *SetupService) ensureProductReady(tile tiles.Definition) error {
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

func (s *SetupService) PoolTillOnline() error {
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

// ERT
func (s *SetupService) ConfigureERT() error {
	return ert.Configure(s.cfg, s.om)
}

func (s *SetupService) UploadERT() error {
	return s.ensureProductReady(ert.Tile)
}

// Service Broker
func (s *SetupService) UploadServiceBroker() error {
	s.pivnet.AcceptEula(service_broker.Tile.Pivnet)
	return s.ensureProductReady(service_broker.Tile)
}

// Stackdriver Nozzle
func (s *SetupService) UploadNozzle() error {
	return s.ensureProductReady(stackdriver_nozzle.Tile)
}
