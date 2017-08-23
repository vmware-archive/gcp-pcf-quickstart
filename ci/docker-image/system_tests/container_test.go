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

package systemTests

import (
	"fmt"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Docker Image", func() {
	commands := map[string][]string{
		"gcloud": []string{"--version"},
		"gsutil": []string{"--version"},
		"ginkgo": []string{"help"},
		"go":     []string{"version"},
		"dig":    []string{"-h"},
	}

	for executable, args := range commands {
		executable := executable
		args := args

		It(fmt.Sprintf("has %v installed", executable), func() {
			command := exec.Command(executable, args...)
			session, err := Start(command, GinkgoWriter, GinkgoWriter)

			Expect(err).ToNot(HaveOccurred())
			Eventually(session, "5s").Should(Exit(0))
		})
	}
})
