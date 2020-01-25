// Copyright © 2019 Thilina Manamgoda
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"github.com/ThilinaManamgoda/password-manager/pkg/storage"
	"github.com/ThilinaManamgoda/password-manager/pkg/utils"
	"github.com/mitchellh/go-homedir"
	"gotest.tools/assert"
	"path/filepath"
	"testing"
)

func TestConfiguration(t *testing.T) {
	err := defaultConf()
	if err != nil {
		t.Error(err)
	}
	result, err := Configuration()
	if err != nil {
		t.Error(err)
	}
	home, err := homedir.Dir()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, filepath.Join(home, "/passwordDB"), result.Storage[storage.ConfKeyPath])
	assert.Equal(t, utils.AESEncryptID, result.EncryptorID)
}
