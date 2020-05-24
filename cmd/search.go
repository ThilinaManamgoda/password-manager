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

// FlagSearchLabel is the flag for searching with labels.
const FlagSearchLabel = "label"

// searchCmd represents the searchId command
var searchCmd = &cobra.Command{
	Use:   "search [ID]",
	Short: "Search Password with ID",
	Long:  `You can use either complete or part of ID/Label for searching`,
	Args:  inputs.HasProvidedValidIDLabel(),
	RunE: func(cmd *cobra.Command, args []string) error {
		mPassword, err := inputs.GetFlagStringVal(cmd, inputs.FlagMasterPassword)
		if err != nil {
			return errors.Wrapf(err, inputs.ErrMsgCannotGetFlag, inputs.FlagMasterPassword)
		}
		if mPassword == "" {
			mPassword, err = inputs.PromptForMPassword()
			if err != nil {
				return errors.Wrap(err, "cannot prompt for Master password")
			}
		}
		showPass, err := inputs.GetFlagBoolVal(cmd, inputs.FlagShowPassword)
		if err != nil {
			return errors.Wrapf(err, inputs.ErrMsgCannotGetFlag, inputs.FlagShowPassword)
		}

		isSearchLabel, err := inputs.GetFlagBoolVal(cmd, FlagSearchLabel)
		if err != nil {
			return errors.Wrapf(err, inputs.ErrMsgCannotGetFlag, FlagSearchLabel)
		}

		conf, err := config.Configuration()
		if err != nil {
			return errors.Wrapf(err, "cannot get configuration")
		}
		passwordRepo, err := passwords.LoadRepo(mPassword, false)
		if err != nil {
			return errors.Wrapf(err, "cannot initialize password repository")
		}

		var passwordEntries []passwords.Entry
		if !isSearchLabel {
			searchID := args[0]
			passwordEntries, err = passwordRepo.SearchEntriesByID(searchID)
			if err != nil {
				return errors.Wrapf(err, "cannot search ID %s", searchID)
			}
		} else {
			label := args[0]
			passwordEntries, err = passwordRepo.SearchLabel(label)
			if err != nil {
				return errors.Wrapf(err, "cannot search Label %s", label)
			}
		}

		if len(passwordEntries) != 0 {
			var psi []inputs.PromptSelectInfo
			for _, e := range passwordEntries {
				psi = append(psi, inputs.PromptSelectInfo{
					ID:          e.ID,
					Description: e.Description,
				})
			}
			sID, err := inputs.PromptForSelect("Choose", conf.SelectListSize, psi)
			if err != nil {
				return errors.Wrap(err, "cannot get prompt for select")
			}
			err = passwordRepo.GetUsernamePassword(sID, showPass)
			if err != nil {
				return errors.Wrapf(err, "cannot get password for ID: %s", sID)
			}
		} else {
			return errors.New("cannot find any match")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().BoolP(inputs.FlagShowPassword, "s", false, "Print password to STDOUT")
	searchCmd.Flags().BoolP(FlagSearchLabel, "l", false, "Search with the Label")
}
