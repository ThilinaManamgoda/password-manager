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
	"github.com/manifoldco/promptui"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
)

// AESEncryptID is the unique identifier for this encryptor
const AESEncryptID = "AES"

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
		EncryptorID:      AESEncryptID,
	}, nil
}

// IsValidByteSlice method check whether the Slice is valid or not
func IsValidByteSlice(data []byte) bool {
	return (data != nil) && (len(data) != 0)
}

// IsPasswordValid method check whether the Password is valid or not
func IsPasswordValid(passphrase string) bool {
	return passphrase != ""
}

// GetFlagStringVal method returns the String flag value
func GetFlagStringVal(cmd *cobra.Command, flag string) (string, error) {
	val, err := cmd.Flags().GetString(flag)
	if err != nil {
		return "", err
	}
	return val, nil
}

// GetFlagBoolVal method returns the Boolean flag value
func GetFlagBoolVal(cmd *cobra.Command, flag string) (bool, error) {
	val, err := cmd.Flags().GetBool(flag)
	if err != nil {
		return false, err
	}
	return val, nil
}

// GetFlagStringArrayVal method returns the String array flag value
func GetFlagStringArrayVal(cmd *cobra.Command, flag string) ([]string, error) {
	val, err := cmd.Flags().GetStringArray(flag)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// IsArgSValid method check whether the CMD args are valid or not
func IsArgSValid(args []string) bool {
	return len(args) != 0
}

// IsArgValid method check whether the CMD arg are valid or not
func IsArgValid(arg string) bool {
	return arg != ""
}

// PromptForString function prompt for string and returns the input
func PromptForString(label string, validate promptui.ValidateFunc) (string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}
	return prompt.Run()
}

// PromptForPassword function prompt for password and returns the input
func PromptForPassword(label string, validate promptui.ValidateFunc)(string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
		Mask:     '*',
	}
	return prompt.Run()
}