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

package inputs

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	// FlagUsername is the flag name for username.
	FlagUsername = "username"
	// FlagPassword is the flag name for password.
	FlagPassword = "password"
	// FlagLabels is the flag name for labels.
	FlagLabels = "labels"
	// FlagMasterPassword is the flag name for masterPassword.
	FlagMasterPassword = "masterPassword"
	// FlagNewMasterPassword is the flag name for new masterPassword.
	FlagNewMasterPassword = "newMasterPassword"
	// FlagShowPassword flag
	FlagShowPassword = "show-pass"
)

// GetFlagStringVal method returns the String flag value.
func GetFlagStringVal(cmd *cobra.Command, flag string) (string, error) {
	return cmd.Flags().GetString(flag)
}

// GetFlagIntVal method returns the int flag value.
func GetFlagIntVal(cmd *cobra.Command, flag string) (int, error) {
	return cmd.Flags().GetInt(flag)
}

// GetFlagBoolVal method returns the Boolean flag value.
func GetFlagBoolVal(cmd *cobra.Command, flag string) (bool, error) {
	return cmd.Flags().GetBool(flag)
}

// GetFlagStringArrayVal method returns the String array flag value.
func GetFlagStringArrayVal(cmd *cobra.Command, flag string) ([]string, error) {
	return cmd.Flags().GetStringArray(flag)
}

// FromFlagsForPasswordEntry functions gets the input values required for Password entry from flags.
func FromFlagsForPasswordEntry(cmd *cobra.Command, uN, password, mPassword *string, labels *[]string) error {
	if uN != nil {
		uNVal, err := GetFlagStringVal(cmd, FlagUsername)
		if err != nil {
			return errors.Wrapf(err, ErrMsgCannotGetFlag, FlagUsername)
		}
		*uN = uNVal
	}

	if password != nil {
		passwordVal, err := GetFlagStringVal(cmd, FlagPassword)
		if err != nil {
			return errors.Wrapf(err, ErrMsgCannotGetFlag, FlagPassword)
		}
		*password = passwordVal
	}

	if labels != nil {
		labelsVal, err := GetFlagStringArrayVal(cmd, FlagLabels)
		if err != nil {
			return errors.Wrapf(err, ErrMsgCannotGetFlag, FlagLabels)
		}
		*labels = labelsVal
	}

	if mPassword != nil {
		mPasswordVal, err := GetFlagStringVal(cmd, FlagMasterPassword)
		if err != nil {
			return errors.Wrapf(err, ErrMsgCannotGetFlag, FlagMasterPassword)
		}
		*mPassword = mPasswordVal
	}
	return nil
}

// IsPasswordValid method check whether the FlagPassword is valid or not.
func IsPasswordValid(passphrase string) bool {
	return passphrase != ""
}
