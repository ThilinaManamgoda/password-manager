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
	"github.com/ThilinaManamgoda/password-manager/pkg/storage"
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
	// EnvPrefix is the environment variable prefix.
	EnvPrefix = "PM"
	// ErrMsgUnableToReadConf is an error message
	ErrMsgUnableToReadConf = "Unable to load configuration file %s"
	// FilePathEnv is env that represents the main configuration path
	FilePathEnv = "PM_CONF_PATH"
	// FlagCSVFile is the CSV file flag.
	FlagCSVFile = "csv-file"
)

// Config struct represent the configuration for the tool.
type Config struct {
	EncryptorID    string  `mapstructure:"encryptorID"`
	SelectListSize int     `mapstructure:"selectListSize"`
	Storage        Storage `mapstructure:"storage"`
}

// TransformedConfig is constructed from Config.
// This configuration is formed to support different kind of components that is been used.
type TransformedConfig struct {
	StorageID      string
	EncryptorID    string
	SelectListSize int
	Storage        map[string]string
}

// Storage represent storage configurations.
type Storage struct {
	File File `mapstructure:"file"`
}

// File represent file storage configurations.
type File struct {
	Path       string `mapstructure:"path"`
	Permission string `mapstructure:"permission"`
}

// Init function configures the viper.
func Init() {
	viper.SetConfigType(YAMLFileType)
	viper.SetEnvPrefix(EnvPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

// Configuration method loads the configuration.
func Configuration() (*TransformedConfig, error) {
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
	parsedConfig := &TransformedConfig{
		EncryptorID:    config.EncryptorID,
		SelectListSize: config.SelectListSize,
	}

	if config.Storage.File.Path != "" {
		storageConf := make(map[string]string)
		storageConf[storage.ConfKeyFilePath] = config.Storage.File.Path
		storageConf[storage.ConfKeyFilePermission] = config.Storage.File.Permission
		parsedConfig.Storage = storageConf
	}
	return parsedConfig, nil
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
	viper.SetDefault("storage.file.path", filepath.Join(home, "/passwordDB"))
	viper.SetDefault("encryptorID", utils.AESEncryptID)
	viper.SetDefault("selectListSize", 5)
	return nil
}
