/*
 *  Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 */

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
