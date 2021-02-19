// Copyright 2021 IBM Corp
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package upload

import (
	"github.com/ppc64le-cloud/pvsadm/test/e2e/framework"
)

const (
	command = "upload"
)

// CMDDescribe annotates the test with the subcommand label.
func CMDDescribe(text string, body func()) bool {
	return framework.Describe(command, text, body)
}
