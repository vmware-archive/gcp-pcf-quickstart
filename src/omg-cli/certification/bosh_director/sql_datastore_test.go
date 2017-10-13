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

package bosh_director_test

import (
	. "omg-cli/certification/environment"

	"omg-cli/ops_manager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SqlDatastore", func() {
	var (
		director *ops_manager.DirectorProperties
	)
	BeforeSuite(func() {
		director = Target().OpsManager().Director()
	})
	// This test fails because we do not use an external SQL database
	XIt("uses an external SQL database", func() {
		// TODO(jrjohnson): This assert will be different but the property to look at will be under .Director
		Expect(director.Director).To(Equal("external"))
	})
})
