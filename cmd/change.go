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
	"github.com/ThilinaManamgoda/password-manager/pkg/config"
	"github.com/ThilinaManamgoda/password-manager/pkg/inputs"
	"github.com/ThilinaManamgoda/password-manager/pkg/passwords"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
)

// changeCmd represents the change command
var changeCmd = &cobra.Command{
	Use:   "change [ID]",
	Short: "Change a password entry",
	Long:  `Change a password entry`,
	Args:  inputs.HasProvidedValidID(),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		mPassword, err := inputs.GetFlagStringVal(cmd, inputs.FlagMasterPassword)
		if err != nil {
			return errors.Wrapf(err, inputs.ErrMsgCannotGetFlag, mPassword)
		}
		if mPassword == "" {
			mPassword, err = inputs.PromptForMPassword()
			if err != nil {
				return errors.Wrap(err, "cannot prompt for Master password")
			}
		}

		conf, err := config.Configuration()
		if err != nil {
			return errors.Wrapf(err, "cannot get configuration")
		}
		passwordRepo, err := passwords.LoadRepo(mPassword, conf.EncryptorID, conf.PasswordDBFilePath)
		if err != nil {
			return errors.Wrapf(err, "cannot initialize password repository")
		}

		passwordEntry, err := passwordRepo.GetPasswordEntry(id)
		if err != nil {
			return errors.Wrapf(err, "cannot get password entry")
		}

		isInteractiveMode, err := inputs.GetFlagBoolVal(cmd, InteractiveMode)
		if err != nil {
			return err
		}
		var uN, password string
		if isInteractiveMode {
			uN, err = inputs.PromptForUsernameWithDefault(passwordEntry.Username)
			if err != nil {
				return errors.Wrap(err, "cannot prompt for username")
			}
			password, err = inputs.PromptForUserPasswordWithDefault(passwordEntry.Password)
			if err != nil {
				return errors.Wrap(err, "cannot prompt for password")
			}
			password,err = inputs.PromptForPasswordSecondTime(password)
			if err != nil {
				return errors.Wrap(err, "cannot prompt for password for the second time")
			}
		} else {
			err = inputs.FromFlagsForPasswordEntry(cmd, &uN, &password, nil, nil)
			if err != nil {
				return errors.Wrapf(err, inputs.ErrMsgCannotGetInput)
			}
		}
		newEntry := passwords.Entry{
			ID:       id,
			Username: uN,
			Password: password,
		}
		err = passwordRepo.ChangePasswordEntry(id, newEntry)
		if err != nil {
			return errors.Wrapf(err, "cannot change password")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(changeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// changeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// changeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	changeCmd.Flags().StringP(inputs.FlagPassword, "p", "", "Password")
	changeCmd.Flags().StringP(inputs.FlagUsername, "u", "", "User Name")
	changeCmd.Flags().BoolP(InteractiveMode, "i", false, "Enable interactive mode")

}
