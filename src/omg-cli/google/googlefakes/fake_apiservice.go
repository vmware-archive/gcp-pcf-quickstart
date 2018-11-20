// Code generated by counterfeiter. DO NOT EDIT.
package googlefakes

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
import (
	"omg-cli/google"
	"sync"
)

type FakeAPIService struct {
	EnableStub        func([]google.API) ([]google.API, error)
	enableMutex       sync.RWMutex
	enableArgsForCall []struct {
		arg1 []google.API
	}
	enableReturns struct {
		result1 []google.API
		result2 error
	}
	enableReturnsOnCall map[int]struct {
		result1 []google.API
		result2 error
	}
	invocations map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeAPIService) Enable(arg1 []google.API) ([]google.API, error) {
	var arg1Copy []google.API
	if arg1 != nil {
		arg1Copy = make([]google.API, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.enableMutex.Lock()
	ret, specificReturn := fake.enableReturnsOnCall[len(fake.enableArgsForCall)]
	fake.enableArgsForCall = append(fake.enableArgsForCall, struct {
		arg1 []google.API
	}{arg1Copy})
	fake.recordInvocation("Enable", []interface{}{arg1Copy})
	fake.enableMutex.Unlock()
	if fake.EnableStub != nil {
		return fake.EnableStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.enableReturns.result1, fake.enableReturns.result2
}

func (fake *FakeAPIService) EnableCallCount() int {
	fake.enableMutex.RLock()
	defer fake.enableMutex.RUnlock()
	return len(fake.enableArgsForCall)
}

func (fake *FakeAPIService) EnableArgsForCall(i int) []google.API {
	fake.enableMutex.RLock()
	defer fake.enableMutex.RUnlock()
	return fake.enableArgsForCall[i].arg1
}

func (fake *FakeAPIService) EnableReturns(result1 []google.API, result2 error) {
	fake.EnableStub = nil
	fake.enableReturns = struct {
		result1 []google.API
		result2 error
	}{result1, result2}
}

func (fake *FakeAPIService) EnableReturnsOnCall(i int, result1 []google.API, result2 error) {
	fake.EnableStub = nil
	if fake.enableReturnsOnCall == nil {
		fake.enableReturnsOnCall = make(map[int]struct {
			result1 []google.API
			result2 error
		})
	}
	fake.enableReturnsOnCall[i] = struct {
		result1 []google.API
		result2 error
	}{result1, result2}
}

func (fake *FakeAPIService) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.enableMutex.RLock()
	defer fake.enableMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeAPIService) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ google.APIService = new(FakeAPIService)
