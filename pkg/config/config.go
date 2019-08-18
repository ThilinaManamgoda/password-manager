// Copyright © 2019 Thilina Manamgoda
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this fileio except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


// Package config holds the functionality for configuration
package config

import (
	"github.com/ThilinaManamgoda/password-manager/pkg/utils"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"path/filepath"
)

// Config struct represent the configuration for the tool
type Config struct {
	PasswordFilePath string `mapstructure:"passwordFilePath"`
	EncryptorID      string `mapstructure:"encryptorID"`
}

// Configuration method loads the configuration
func Configuration() (*Config, error) {
	var config, err = defaultConf()
	if err != nil {
		return nil, errors.Wrap(err, "cannot load default config")
	}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal the configuration")
	}
	return config, nil
}

func defaultConf() (*Config, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, errors.Wrap(err, "cannot retrieve Home directory path")
	}
	return &Config{
		PasswordFilePath: filepath.Join(home, "/passwordDB"),
		EncryptorID:      utils.AESEncryptID,
	}, nil
}

