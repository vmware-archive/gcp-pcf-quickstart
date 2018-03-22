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

package google

type ZoneResult struct {
	Region string
	Zone1  string
	Zone2  string
	Zone3  string
}

func ParseZone(zone string) (*ZoneResult, error) {
	// TODO: Implement a real version
	return &ZoneResult{Region: "us-east1", Zone1: "us-east1-b", Zone2: "us-east1-c", Zone3: "us-east1-d"}, nil
}
