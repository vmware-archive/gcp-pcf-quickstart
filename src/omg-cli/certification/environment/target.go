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

package environment

import (
	"net/http"

	"sync"

	"fmt"
	"os"

	"omg-cli/config"

	"omg-cli/ops_manager"

	"log"

	"github.com/onsi/ginkgo"
)

type TargetSite interface {
	// OpsManager returns an OpsManagerQuery that can be used
	// to access properties about the target Ops Manager
	OpsManager() OpsManagerQuery
	// GoogleClient returns an authenticated http client that
	// can be used to create service objects for GCP APIs
	GoogleClient() *http.Client
}

var target TargetSite
var targetOnce sync.Once

const envDirName = "ENV_DIR"

func Target() TargetSite {
	targetOnce.Do(func() {
		envDir := os.Getenv(envDirName)
		if envDir == "" {
			ginkgo.Fail(fmt.Sprintf("missing test data, expected environment variable %s to contain path", envDirName))
		}

		cfg, err := config.FromTerraformDirectory(envDir)
		if err != nil {
			ginkgo.Fail(fmt.Sprintf("loading terraform state: %v", err))
		}

		target = &liveTarget{cfg: cfg}
	})
	return target
}

type liveTarget struct {
	cfg *config.Config
}

func (lt *liveTarget) OpsManager() OpsManagerQuery {
	logger := log.New(os.Stdout, "TODO(jrjohnson): test logger", 0)
	omSdk, err := ops_manager.NewSdk(fmt.Sprintf("https://%s", lt.cfg.OpsManagerHostname), lt.cfg.OpsManager, *logger)
	if err != nil {
		ginkgo.Fail(fmt.Sprintf("creating ops manager sdk: %v", err))
	}

	return &liveOpsManager{sdk: omSdk}
}

func (lt *liveTarget) GoogleClient() *http.Client {
	panic("implement me")
}
