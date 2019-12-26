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

// Package config holds the functionality for configuration.
package config

import (
	"github.com/ThilinaManamgoda/password-manager/pkg/utils"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

const (
	// YAMLFileType is the file type YAML.
	YAMLFileType = "yaml"
	// EnvPrefix is the environment variable prefix
	EnvPrefix = "PM"
	// ErrMsgUnableToReadConf is an error message
	ErrMsgUnableToReadConf = "Unable to load configuration file %s"
	// FilePathEnv is env that represents the main configuration path
	FilePathEnv = "PM_CONF_PATH"
)

// Config struct represent the configuration for the tool.
type Config struct {
	PasswordDBFilePath string `mapstructure:"passwordDBFilePath"`
	EncryptorID        string `mapstructure:"encryptorID"`
}

// Init function configures the viper.
func Init() {
	viper.SetConfigType(YAMLFileType)
	viper.SetEnvPrefix(EnvPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

// Configuration method loads the configuration.
func Configuration() (*Config, error) {
	err := defaultConf()
	if err != nil {
		return nil, errors.Wrap(err, "cannot set default config")
	}

	err = loadConfigFile()
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal the configuration")
	}
	return config, nil
}

func loadConfigFile() error {
	confFile, exists := os.LookupEnv(FilePathEnv)
	if exists {
		viper.SetConfigFile(confFile)
		if err := viper.ReadInConfig(); err != nil {
			return errors.Wrapf(err, ErrMsgUnableToReadConf, confFile)
		}
	}
	return nil
}

func defaultConf() error {
	home, err := homedir.Dir()
	if err != nil {
		return errors.Wrap(err, "cannot retrieve Home directory path")
	}
	viper.SetDefault("passwordDBFilePath", filepath.Join(home, "/passwordDB"))
	viper.SetDefault("encryptorID", utils.AESEncryptID)
	return nil
}
