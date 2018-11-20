package ops_manager

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

// This file contains types used for marshalling requests/responses to Ops Manager.
// These types are public so they are visible to the `json` package. They should
// not be consumed by other packages.

type CredentialResponse struct {
	Credential CredentialWrapper `json:"credential"`
}

type CredentialWrapper struct {
	Type  string           `json:"type"`
	Value SimpleCredential `json:"value"`
}

type Product struct {
	InstallationName string `json:"installation_name"`
	GUID             string `json:"guid"`
	Type             string `json:"type"`
}

type ErrorResponse struct {
	Errors map[string][]string `json:errors`
}

type JobsResponse struct {
	Jobs []Job `json:"jobs"`
}

type UnlockRequest struct {
	Passphrase string `json:"passphrase"`
}
