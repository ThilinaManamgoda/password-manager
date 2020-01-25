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
	"github.com/ThilinaManamgoda/password-manager/pkg/inputs"
	"github.com/ThilinaManamgoda/password-manager/pkg/passwords"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// changeMasterPasswordCmd represents the changeMasterPassword command
var changeMasterPasswordCmd = &cobra.Command{
	Use:   "change-master-password",
	Short: "Change Master password",
	Long:  `Change Master password`,
	RunE: func(cmd *cobra.Command, args []string) error {
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
		newMasterPassword, err := inputs.GetFlagStringVal(cmd, inputs.FlagNewMasterPassword)
		if err != nil {
			return errors.Wrapf(err, inputs.ErrMsgCannotGetFlag, mPassword)
		}
		if newMasterPassword == "" {
			newMasterPassword, err = inputs.PromptForNewMPassword()
			if err != nil {
				return errors.Wrap(err, "cannot prompt for new password")
			}
			newMasterPassword, err = inputs.PromptForMPasswordSecondTime(newMasterPassword)
			if err != nil {
				return errors.Wrap(err, "cannot prompt for new password again")
			}

		}

		passwordRepo, err := passwords.LoadRepo(mPassword)
		if err != nil {
			return errors.Wrapf(err, "cannot initialize password repository")
		}
		err = passwordRepo.ChangeMasterPassword(newMasterPassword)
		if err != nil {
			return errors.Wrapf(err, "cannot change Master Password")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(changeMasterPasswordCmd)
	changeMasterPasswordCmd.Flags().StringP(inputs.FlagNewMasterPassword, "n", "", "New Master password")
}
