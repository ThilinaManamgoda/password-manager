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

package cmd

import (
	"github.com/password-manager/pkg/encrypt"
	"github.com/password-manager/pkg/passwords"
	"github.com/password-manager/pkg/utils"
	"github.com/spf13/cobra"
)

const (
	Username       = "username"
	Password       = "password"
	Labels         = "labels"
	MasterPassword = "masterPassword"
	Id             = "id"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new password",
	Long:  `Add a new password`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := utils.GetFlagStringVal(cmd, Id)
		if err != nil {
			return err
		}
		uN, err := utils.GetFlagStringVal(cmd,Username)
		if err != nil {
			return err
		}
		password, err := utils.GetFlagStringVal(cmd,Password)
		if err != nil {
			return err
		}
		labels, err := utils.GetFlagStringArrayVal(cmd,Labels)
		if err != nil {
			return err
		}
		mPassword, err := utils.GetFlagStringVal(cmd,MasterPassword)
		if err != nil {
			return err
		}
		config, err := utils.Configuration()
		if err != nil {
			return err
		}
		encriptorFac := &encrypt.Factory{
			ID: config.EncryptorID,
		}
		passwordRepo := &passwords.PasswordRepository{
			MasterPassword: mPassword,
			Encryptor:      encriptorFac.GetEncryptor(),
			PasswordFile: utils.PasswordFile{
				File: config.PasswordFilePath,
			},
		}
		err = passwordRepo.Add(id, uN, password, labels)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Pexxrsistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	addCmd.Flags().StringP(Password, "p", "", "Password")
	addCmd.Flags().StringP(Username, "u", "", "User Name")
	addCmd.Flags().StringP(Id, "i", "", "Id for entry")
	addCmd.Flags().StringArrayP(Labels, "l", nil, "Labels for the entry")
}
