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

package cmd

import (
	"github.com/ThilinaManamgoda/password-manager/pkg/config"
	"github.com/ThilinaManamgoda/password-manager/pkg/inputs"
	"github.com/ThilinaManamgoda/password-manager/pkg/passwords"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	// ID flag
	ID = "id"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [ID]",
	Short: "Add a new password",
	Long:  `Add a new password`,
	Args:  inputs.HasProvidedValidID(),
	RunE: func(cmd *cobra.Command, args []string) error {
		isInteractiveMode, err := inputs.GetFlagBoolVal(cmd, InteractiveMode)
		if err != nil {
			return err
		}
		var uN, password, mPassword string
		var labels []string
		if isInteractiveMode {
			err := inputs.FromPromptForPasswordEntry(&uN, &password, &mPassword, &labels)
			if err != nil {
				return errors.Wrapf(err, inputs.ErrMsgCannotGetInput)
			}
		} else {
			err := inputs.FromFlagsForPasswordEntry(cmd, &uN, &password, &mPassword, &labels)
			if err != nil {
				return errors.Wrapf(err, inputs.ErrMsgCannotGetInput)
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

		id := args[0]
		err = passwordRepo.Add(id, uN, password, labels)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringP(inputs.FlagPassword, "p", "", "Password")
	addCmd.Flags().StringP(inputs.FlagUsername, "u", "", "User Name")
	addCmd.Flags().StringArrayP(inputs.FlagLabels, "l", nil, "Labels for the password entry")
	addCmd.Flags().BoolP(InteractiveMode, "i", false, "Enable interactive mode")
}
