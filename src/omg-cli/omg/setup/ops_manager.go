/*
 * Copyright 2017 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package setup

import (
	"omg-cli/config"
	"omg-cli/ops_manager"
	"omg-cli/pivnet"

	"errors"
	"fmt"
	"log"
	"omg-cli/omg/tiles"
	"os"
	"path"
	"time"
)

type OpsManager struct {
	cfg    *config.Config
	om     *ops_manager.Sdk
	pivnet *pivnet.Sdk
	logger *log.Logger
	tiles  []tiles.TileInstaller
}

func NewService(cfg *config.Config, omSdk *ops_manager.Sdk, pivnetSdk *pivnet.Sdk, logger *log.Logger, tiles []tiles.TileInstaller) *OpsManager {
	return &OpsManager{cfg, omSdk, pivnetSdk, logger, tiles}
}

func (s *OpsManager) SetupAuth() error {
	return s.om.SetupAuth()
}

func (s *OpsManager) Unlock() error {
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

func (s *OpsManager) ApplyChanges() error {
	return s.om.ApplyChanges()
}

func (s *OpsManager) productInstalled(product config.OpsManagerMetadata) (bool, error) {
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

func (s *OpsManager) ensureProductReady(tile config.Tile) error {
	if i, err := s.productInstalled(tile.Product); i == true || err != nil {
		return err
	}

	return s.uploadProduct(tile.Pivnet)
}

func (s *OpsManager) uploadProduct(tile config.PivnetMetadata) error {
	file, err := s.pivnet.DownloadTile(tile)
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	return s.om.UploadProduct(file.Name())
}

func (s *OpsManager) uploadStemcell(tile config.StemcellMetadata) error {
	file, err := s.pivnet.DownloadTile(tile.PivnetMetadata)
	if err != nil {
		return err
	}

	newPath := path.Join(path.Dir(file.Name()), fmt.Sprintf("%s.tgz", tile.StemcellName))
	if err := os.Rename(file.Name(), newPath); err != nil {
		os.Remove(file.Name())
		return fmt.Errorf("unable to rename download stemcell: %v", err)
	}
	defer os.Remove(newPath)

	return s.om.UploadStemcell(newPath)
}

func (s *OpsManager) PoolTillOnline() error {
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

func (s *OpsManager) ConfigureTiles() error {
	for _, t := range s.tiles {
		s.logger.Printf("configuring tile: %s", t.Definition().Product.Name)
		if err := t.Configure(s.cfg, s.om); err != nil {
			return err
		}
	}

	return nil
}

func (s *OpsManager) UploadTiles() error {
	for _, t := range s.tiles {
		if t.BuiltIn() {
			continue
		}

		if s.cfg.PivnetAcceptEula {
			if err := s.pivnet.AcceptEula(t.Definition().Pivnet); err != nil {
				s.logger.Printf("warning, unable to accept eula: %v", err)
			}
		}

		if err := s.ensureProductReady(t.Definition()); err != nil {
			return err
		}

		if stemcell := t.Definition().Stemcell; stemcell != nil {
			if err := s.uploadStemcell(*stemcell); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *OpsManager) DeleteInstallation() error {
	return s.om.DeleteInstallation()
}
