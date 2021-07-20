// +build !linux

/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package enricher

import "errors"

var errUnsupoprtedPlatform = errors.New("unsupoprted platform")

// syscallName returns the syscall name for the provided ID.
func syscallName(id int32) (string, error) {
	return "", errUnsupoprtedPlatform
}
