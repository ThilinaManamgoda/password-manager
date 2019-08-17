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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	// Username flag
	Username       = "username"
	// Password flag
	Password       = "password"
	// Labels flag
	Labels         = "labels"
	// MasterPassword flag
	MasterPassword = "masterPassword"
	// ID flag
	ID             = "id"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [ID]",
	Short: "Add a new password",
	Long:  `Add a new password`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if ! utils.IsArgSValid(args) {
			return errors.New("invalid argument")
		}
		if ! utils.IsArgValid(args[0]) {
			return errors.New("invalid ID")
		}
		id := args[0]

		isInteractiveMode, err := utils.GetFlagBoolVal(cmd, InteractiveMode)
		if err != nil {
			return err
		}
		var uN, password, mPassword string
		var labels []string
		if isInteractiveMode {
			uN, err = promptForUsername()
			if err != nil {
				return errors.Wrap(err, "cannot prompt for Username")
			}
			password, err = promptForPassword()
			if err != nil {
				return errors.Wrap(err, "cannot prompt for password")
			}
			mPassword, err = promptForMPassword()
			if err != nil {
				return errors.Wrap(err, "cannot prompt for Master password")
			}
		} else {
			uN, err = utils.GetFlagStringVal(cmd, Username)
			if err != nil {
				return errors.Wrapf(err, "cannot get value of %s flag", Username)
			}
			password, err = utils.GetFlagStringVal(cmd, Password)
			if err != nil {
				return errors.Wrapf(err, "cannot get value of %s flag", Password)
			}
			labels, err = utils.GetFlagStringArrayVal(cmd, Labels)
			if err != nil {
				return errors.Wrapf(err, "cannot get value of %s flag", Labels)
			}
			mPassword, err = utils.GetFlagStringVal(cmd, MasterPassword)
			if err != nil {
				return errors.Wrapf(err, "cannot get value of %s flag", MasterPassword)
			}
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

func promptForUsername()(string, error) {
	validate := func(input string) error {
		if len(input) < 3 {
			return errors.New("username must have more than 3 characters")
		}
		return nil
	}
	return utils.PromptForString("Username", validate)
}

func promptForMPassword()(string, error) {
	validate := func(input string) error {
		if len(input) < 6 {
			return errors.New("master password must have more than 6 characters")
		}
		return nil
	}
	return utils.PromptForPassword("Master password", validate)
}

func promptForPassword()(string, error) {
	validate := func(input string) error {
		if len(input) < 6 {
			return errors.New("password must have more than 6 characters")
		}
		return nil
	}
	return utils.PromptForPassword("Password", validate)
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
	addCmd.Flags().StringArrayP(Labels, "l", nil, "Labels for the entry")
}
