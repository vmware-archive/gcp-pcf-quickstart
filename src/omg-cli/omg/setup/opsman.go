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
	"errors"
	"fmt"
	"log"
	"omg-cli/omg/tiles"
	"os"
	"path"
	"time"

	"omg-cli/config"
	"omg-cli/opsman"
	"omg-cli/pivnet"
)

// OpsManager is used for setting up a real Ops Manager.
type OpsManager struct {
	cfg       *config.Config
	envCfg    *config.EnvConfig
	om        *opsman.Sdk
	pivnet    *pivnet.Sdk
	logger    *log.Logger
	tiles     []tiles.TileInstaller
	tileCache *pivnet.TileCache
}

// NewOpsManager creates a new OpsManager for setup purposes.
func NewOpsManager(cfg *config.Config, envCfg *config.EnvConfig, omSdk *opsman.Sdk, pivnetSdk *pivnet.Sdk, logger *log.Logger, tiles []tiles.TileInstaller, tileCache *pivnet.TileCache) *OpsManager {
	return &OpsManager{cfg, envCfg, omSdk, pivnetSdk, logger, tiles, tileCache}
}

// SetupAuth configures the initial username, password, and decryptionPhrase
func (s *OpsManager) SetupAuth() error {
	return s.om.SetupAuth()
}

// Unlock decrypts Ops Manager. This is needed after a reboot before attempting to authenticate.
// This task runs asynchronously. Query the status by invoking ReadyForAuth.
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
			return errors.New("timeout waiting for Ops Manager to unlock")
		case <-timer:
			if s.om.ReadyForAuth() {
				return nil
			}
			s.logger.Print("waiting for Ops Manager to unlock")
			timer = time.After(time.Duration(5 * time.Second))
		}
	}
}

// ApplyChangesPAS runs apply_changes on all tiles except for ones which depend on the PAS.
func (s *OpsManager) ApplyChangesPAS() error {
	var args []string
	for _, tile := range s.tiles {
		if !tile.Definition(s.envCfg).Product.DependsOnPAS {
			name := tile.Definition(s.envCfg).Product.Name
			args = append(args, fmt.Sprintf("--product-name=%s", name))
		}
	}
	return s.om.ApplyChanges(args)
}

// ApplyChangesSkipUnchanged runs apply_changes on any tiles which are not up to date.
func (s *OpsManager) ApplyChangesSkipUnchanged() error {
	return s.om.ApplyChanges([]string{"--skip-unchanged-products"})
}

// ApplyDirector runs apply_changes on the BOSH Director.
func (s *OpsManager) ApplyDirector() error {
	return s.om.ApplyDirector()
}

func (s *OpsManager) productInstalled(product config.OpsManagerMetadata) (bool, error) {
	products, err := s.om.AvailableProducts()
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
	if i, err := s.productInstalled(tile.Product); i || err != nil {
		return err
	}

	return s.uploadProduct(tile.Pivnet)
}

func (s *OpsManager) uploadProduct(tile config.PivnetMetadata) error {
	file, err := s.tileCache.Open(tile)

	if file == nil {
		s.logger.Printf("tile not found in cache, downloading")
		file, err = s.pivnet.DownloadTile(tile)
		defer os.Remove(file.Name())
	}

	if err != nil {
		return err
	}

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

// PoolTillOnline waits for the Ops Manager's API to be available.
func (s *OpsManager) PoolTillOnline() error {
	timer := time.After(time.Duration(0 * time.Second))
	timeout := time.After(time.Duration(240 * time.Second))
	for {
		select {
		case <-timeout:
			return errors.New("timeout waiting for Ops Manager to start")
		case <-timer:
			if s.om.Online() {
				return nil
			}
			s.logger.Print("waiting for Ops Manager to start")
			timer = time.After(time.Duration(5 * time.Second))
		}
	}
}

// ConfigureTiles configures each tile.
func (s *OpsManager) ConfigureTiles() error {
	for _, t := range s.tiles {
		s.logger.Printf("configuring tile: %s", t.Definition(s.envCfg).Product.Name)
		if err := t.Configure(s.envCfg, s.cfg, s.om); err != nil {
			return err
		}
	}

	return nil
}

// UploadTiles uploads the tiles and their stemcells to the Ops Manager.
func (s *OpsManager) UploadTiles() error {
	for _, t := range s.tiles {
		if t.BuiltIn() {
			continue
		}

		if err := s.ensureProductReady(t.Definition(s.envCfg)); err != nil {
			return err
		}

		if stemcell := t.Definition(s.envCfg).Stemcell; stemcell != nil {
			if err := s.uploadStemcell(*stemcell); err != nil {
				return err
			}
		}
	}

	return nil
}

// DeleteInstallation deletes an installation.
func (s *OpsManager) DeleteInstallation() error {
	return s.om.DeleteInstallation()
}
