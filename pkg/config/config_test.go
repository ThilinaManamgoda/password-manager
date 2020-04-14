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
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfiguration(t *testing.T) {
	Init()
	result, err := Configuration()
	if err != nil {
		t.Error(err)
	}
	home, err := homedir.Dir()
	if err != nil {
		t.Error(err)
	}
	directoryPath := filepath.Join(home, "/"+DefaultDirectoryName)
	assert.Equal(t, directoryPath, result.DirectoryPath)
	assert.Equal(t, filepath.Join(directoryPath, "/"+DefaultPasswordDBFile), result.Storage[storage.ConfKeyFilePath])
	assert.Equal(t, DefaultFilePermission, result.Storage[storage.ConfKeyFilePermission])
	assert.Equal(t, utils.AESEncryptID, result.EncryptorID)
	assert.Equal(t, DefaultSelectListSize, result.SelectListSize)
}

func TestConfigurationWithEnv(t *testing.T) {
	Init()
	err := os.Setenv("PM_STORAGE_GOOGLEDRIVE_ENABLE", "true")
	if err != nil {
		t.Error(err)
	}
	tmpPath := "/root/user/test"
	err = os.Setenv("PM_DIRECTORYPATH", tmpPath)
	if err != nil {
		t.Error(err)
	}
	result, err := Configuration()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, tmpPath, result.DirectoryPath)
	assert.Equal(t, DefaultPasswordDBFile, result.Storage[storage.ConfKeyPasswordDBFile])
	assert.Equal(t, filepath.Join(tmpPath, "/"+DefaultTokenFileName), result.Storage[storage.ConfKeyTokenFilePath])
	assert.Equal(t, DefaultDirectoryName, result.Storage[storage.ConfKeyDirectory])
}
