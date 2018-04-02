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

import (
	"fmt"
	"strings"

	compute "google.golang.org/api/compute/v1"
)

type ZoneResult struct {
	Region string
	Zone1  string
	Zone2  string
	Zone3  string
}

// Get the last path in a URL
//
// eg: https://www.googleapis.com/compute/v1/projects/<project>/regions/us-central1 => us-central1
func suffix(url string) string {
	urlParts := strings.Split(url, "/")
	return urlParts[len(urlParts)-1]
}

func ParseZone(project, baseZone string, client *compute.Service) (res ZoneResult, err error) {
	zone, err := client.Zones.Get(project, baseZone).Do()
	if err != nil {
		return
	}
	res.Region = suffix(zone.Region)

	region, err := client.Regions.Get(project, res.Region).Do()
	if err != nil {
		return
	}
	if zones := len(region.Zones); zones < 3 {
		err = fmt.Errorf("region %s does not contain enough zones, found: %d", region.Name, zones)
		return
	}

	res.Zone1 = suffix(region.Zones[0])
	res.Zone2 = suffix(region.Zones[1])
	res.Zone3 = suffix(region.Zones[2])

	return
}
