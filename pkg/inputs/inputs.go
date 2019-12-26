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

// Package inputs handles the user interactions
package inputs

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strings"
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
	// PromptUsername is the prompt name for username
	PromptUsername = "Username"
	// PromptPassword is the prompt name for password
	PromptPassword = "Password"
	// PromptLabels is the flag name for labels.
	PromptLabels = "Labels"
	// PromptMasterPassword is the flag name for masterPassword.
	PromptMasterPassword = "Master password"
	// ErrMSGCannotPrompt is an error message
	ErrMSGCannotPrompt = "cannot prompt for %s"
	// ErrMsgCannotGetInput is an error message
	ErrMsgCannotGetInput = "cannot get input"
	// ErrMsgCannotGetFlag is an error message
	ErrMsgCannotGetFlag = "cannot get value of %s flag"
	// ErrMsgMasterPasswordMissMatch is an error message
	ErrMsgMasterPasswordMissMatch = "master password doesn't match"
	// ErrMsgPasswordMissMatch is an error message
	ErrMsgPasswordMissMatch = "password doesn't match"
	// MinPasswordCharacters is the allowed minimum number of characters for a password.
	MinPasswordCharacters = 6
)

var (
	userNameValidator = func(input string) error {
		if len(input) < MinPasswordCharacters {
			return errors.New(fmt.Sprintf("username must have more than %d characters", MinPasswordCharacters))
		}
		return nil
	}

	passwordValidator = func(errorMsg string) func(input string) error {
		return func(input string) error {
			if len(input) < MinPasswordCharacters {
				return errors.New(errorMsg)
			}
			return nil
		}
	}

	passwordEqualityValidator = func(currentPassword, errorMsg string) func(input string) error {
		return func(input string) error {
			if currentPassword != input {
				return errors.New(errorMsg)
			}
			return nil
		}
	}

	masterPasswordValidator    = passwordValidator(fmt.Sprintf("master password must have more than %d characters", MinPasswordCharacters))
	userPasswordValidator      = passwordValidator(fmt.Sprintf("password must have more than %d characters", MinPasswordCharacters))
	newMasterPasswordValidator = passwordValidator(fmt.Sprintf("new master password must have more than %d characters", MinPasswordCharacters))
)

// IsPasswordValid method check whether the FlagPassword is valid or not.
func IsPasswordValid(passphrase string) bool {
	return passphrase != ""
}

// GetFlagStringVal method returns the String flag value.
func GetFlagStringVal(cmd *cobra.Command, flag string) (string, error) {
	val, err := cmd.Flags().GetString(flag)
	if err != nil {
		return "", err
	}
	return val, nil
}

// GetFlagIntVal method returns the int flag value.
func GetFlagIntVal(cmd *cobra.Command, flag string) (int, error) {
	val, err := cmd.Flags().GetInt(flag)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// GetFlagBoolVal method returns the Boolean flag value.
func GetFlagBoolVal(cmd *cobra.Command, flag string) (bool, error) {
	val, err := cmd.Flags().GetBool(flag)
	if err != nil {
		return false, err
	}
	return val, nil
}

// GetFlagStringArrayVal method returns the String array flag value.
func GetFlagStringArrayVal(cmd *cobra.Command, flag string) ([]string, error) {
	val, err := cmd.Flags().GetStringArray(flag)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// IsValidSingleArg method check whether the CMD args are valid or not.
func IsValidSingleArg(args []string) bool {
	return len(args) == 1
}

// IsArgValid method check whether the CMD arg are valid or not.
func IsArgValid(arg string) bool {
	return arg != ""
}

// PromptForString function prompt for string and returns the input.
func PromptForString(label string, validate promptui.ValidateFunc) (string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}
	return prompt.Run()
}

// PromptForStringWithDefault function prompts for a string with a choice of default value and returns the chosen value.
func PromptForStringWithDefault(label, defaultVal string, validate promptui.ValidateFunc) (string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
		Default:  defaultVal,
	}
	return prompt.Run()
}

// PromptForUserPasswordWithDefault function prompts for a user password with a choice of default value and returns the chosen value.
func PromptForUserPasswordWithDefault(defaultVal string) (string, error) {
	return PromptForPasswordWithDefault(PromptPassword, defaultVal, userPasswordValidator)
}

// PromptForPasswordWithDefault function prompts for a password with a choice of default value and returns the chosen value.
func PromptForPasswordWithDefault(label, defaultVal string, validate promptui.ValidateFunc) (string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
		Default:  defaultVal,
		Mask:     '*',
	}
	return prompt.Run()
}

func promptForPassword(label string, validate promptui.ValidateFunc) (string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
		Mask:     '*',
	}
	return prompt.Run()
}

// PromptForSelect start selection and return the selected value.
func PromptForSelect(l string, items []string) (string, error) {
	prompt := promptui.Select{
		Label: l,
		Items: items,
	}
	_, result, err := prompt.Run()
	return result, err
}

// HasProvidedValidID returns a function which validates the ID input.
func HasProvidedValidID() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if !IsValidSingleArg(args) {
			return errors.New("please give a valid ID")
		}
		if !IsArgValid(args[0]) {
			return errors.New("id cannot be empty")
		}
		return nil
	}
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

// FromPromptForPasswordEntry functions gets the input values required for Password entry by prompting.
func FromPromptForPasswordEntry(uN, password, mPassword *string, labels *[]string) error {
	uNVal, err := PromptForUsername()
	if err != nil {
		return errors.Wrapf(err, ErrMSGCannotPrompt, PromptUsername)
	}
	*uN = uNVal
	passwordVal, err := PromptForPassword()
	if err != nil {
		return errors.Wrapf(err, ErrMSGCannotPrompt, PromptPassword)
	}
	*password = passwordVal

	passwordVal, err = PromptForPasswordSecondTime(passwordVal)
	if err != nil {
		return errors.Wrapf(err, ErrMSGCannotPrompt, "Again password")
	}

	labelsVal, err := PromptForLabels()
	if err != nil {
		return errors.Wrapf(err, ErrMSGCannotPrompt, PromptLabels)
	}
	*labels = labelsVal
	mPasswordVal, err := PromptForMPassword()
	if err != nil {
		return errors.Wrapf(err, ErrMSGCannotPrompt, PromptMasterPassword)
	}
	*mPassword = mPasswordVal
	return nil
}

// PromptForUsername prompt for username and returns the value.
func PromptForUsername() (string, error) {
	return PromptForString(PromptUsername, userNameValidator)
}

// PromptForUsernameWithDefault prompt for username with a deault value and returns the chosen value.
func PromptForUsernameWithDefault(defaultVal string) (string, error) {
	return PromptForStringWithDefault(PromptUsername, defaultVal, userNameValidator)
}

// PromptForLabels prompts for labels and returns the given labels.
func PromptForLabels() ([]string, error) {
	validate := func(input string) error {
		return nil
	}
	lInput, err := PromptForString(PromptLabels, validate)
	if err != nil {
		return nil, nil
	}
	l := strings.Split(lInput, ",")
	if len(l) == 0 {
		return nil, nil
	}
	return l, nil
}

// PromptForNewMPassword prompts for a new master password and returns master password.
func PromptForNewMPassword() (string, error) {
	return promptForPassword("New Master password", newMasterPasswordValidator)
}

// PromptForMPassword prompts for the master password and returns master password.
func PromptForMPassword() (string, error) {
	return promptForPassword(PromptMasterPassword, masterPasswordValidator)
}

// PromptForPassword function prompts for password and returns the input.
func PromptForPassword() (string, error) {
	return promptForPassword(PromptPassword, userPasswordValidator)
}

// PromptForPasswordSecondTime function prompts for password for second time to validate and returns the input.
func PromptForPasswordSecondTime(currentPassword string) (string, error) {
	return promptForPassword("Enter the Password again", passwordEqualityValidator(currentPassword, ErrMsgPasswordMissMatch))
}

// PromptForMPasswordSecondTime function prompts for master password for second time to validate and returns the input.
func PromptForMPasswordSecondTime(currentPassword string) (string, error) {
	return promptForPassword("Enter the Master password again", passwordEqualityValidator(currentPassword, ErrMsgMasterPasswordMissMatch))
}
