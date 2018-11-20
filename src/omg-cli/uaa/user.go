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

package uaa

type Name struct {
	FamilyName string `json:"familyName"`
	GivenName  string `json:"givenName"`
}
type PhoneNumber struct {
	Value string `json:"value"`
}
type Email struct {
	Value   string `json:"value"`
	Primary bool   `json:"primary"`
}
type User struct {
	ID           string        `json:"id"`
	Username     string        `json:"userName"`
	Password     string        `json:"password"`
	Name         Name          `json:"name"`
	PhoneNumbers []PhoneNumber `json:"phoneNumbers"`
	Emails       []Email       `json:"emails"`
	Active       bool          `json:"active"`
}
