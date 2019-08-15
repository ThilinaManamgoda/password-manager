// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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

package utils

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
)

// AESEncryptID is the unique identifier for this encryptor
const AESEncryptID = "AES"

type Config struct {
	PasswordFilePath string `mapstructure:"passwordFilePath"`
	EncryptorID      string `mapstructure:"encryptorID"`
}

func Configuration() (*Config, error) {
	var config, err = defaultConf()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func defaultConf() (*Config, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	return &Config{
		PasswordFilePath: filepath.Join(home, "/passwordDB"),
		EncryptorID:      AESEncryptID,
	}, nil
}

func IsValidByteSlice(data []byte) bool {
	return (data != nil) && (len(data) != 0)
}

func IsPasswordValid(passphrase string) bool {
	return passphrase != ""
}

func GetFlagStringVal(cmd *cobra.Command, flag string) (string, error) {
	val, err := cmd.Flags().GetString(flag)
	if err != nil {
		return "", err
	}
	return val, nil
}

func GetFlagBoolVal(cmd *cobra.Command, flag string) (bool, error) {
	val, err := cmd.Flags().GetBool(flag)
	if err != nil {
		return false, err
	}
	return val, nil
}

func GetFlagStringArrayVal(cmd *cobra.Command, flag string) ([]string, error) {
	val, err := cmd.Flags().GetStringArray(flag)
	if err != nil {
		return nil, err
	}
	return val, nil
}
