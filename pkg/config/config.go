/*
 *  Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 */

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

