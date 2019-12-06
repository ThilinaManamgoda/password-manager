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
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strings"
)

const (
	Username = "username"

	Password = "password"

	Labels = "labels"

	ErrMSGCannotPrompt = "cannot prompt for %s"

	ErrMsgCannotGetInput = "cannot get input"

	ErrMSGCannotGetFlag = "cannot get value of %s flag"

	MasterPassword = "masterPassword"
)

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

// IsValidSingleArg method check whether the CMD args are valid or not
func IsValidSingleArg(args []string) bool {
	return len(args) == 1
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

func promptForPassword(label string, validate promptui.ValidateFunc) (string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
		Mask:     '*',
	}
	return prompt.Run()
}

// PromptForSelect start selection and return the selected value
func PromptForSelect(l string, items []string) (string, error) {
	prompt := promptui.Select{
		Label: l,
		Items: items,
	}
	_, result, err := prompt.Run()
	return result, err
}

// HasProvidedValidID returns a function which validates the ID input
func HasProvidedValidID() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if !IsValidSingleArg(args) {
			return errors.New("Please give a valid ID")
		}
		if !IsArgValid(args[0]) {
			return errors.New("id cannot be empty")
		}
		return nil
	}
}

func FromFlags(cmd *cobra.Command, uN, password, mPassword *string, labels *[]string) error {
	uNVal, err := GetFlagStringVal(cmd, Username)
	if err != nil {
		return errors.Wrapf(err, ErrMSGCannotGetFlag, Username)
	}
	*uN = uNVal
	passwordVal, err := GetFlagStringVal(cmd, Password)
	if err != nil {
		return errors.Wrapf(err, ErrMSGCannotGetFlag, Password)
	}
	*password = passwordVal
	labelsVal, err := GetFlagStringArrayVal(cmd, Labels)
	if err != nil {
		return errors.Wrapf(err, ErrMSGCannotGetFlag, Labels)
	}
	*labels = labelsVal
	mPasswordVal, err := GetFlagStringVal(cmd, MasterPassword)
	if err != nil {
		return errors.Wrapf(err, ErrMSGCannotGetFlag, MasterPassword)
	}
	*mPassword = mPasswordVal
	return nil
}

func FromPrompt(uN, password, mPassword *string, labels *[]string) error {
	uNVal, err := PromptForUsername()
	if err != nil {
		return errors.Wrapf(err, ErrMSGCannotPrompt, "Username")
	}
	*uN = uNVal
	passwordVal, err := PromptForPassword()
	if err != nil {
		return errors.Wrapf(err, ErrMSGCannotPrompt, "Password")
	}
	*password = passwordVal
	labelsVal, err := PromptForLabels()
	if err != nil {
		return errors.Wrapf(err, ErrMSGCannotPrompt, "Labels")
	}
	*labels = labelsVal
	mPasswordVal, err := PromptForMPassword()
	if err != nil {
		return errors.Wrapf(err, ErrMSGCannotPrompt, "Master password")
	}
	*mPassword = mPasswordVal
	return nil
}

func PromptForUsername() (string, error) {
	validate := func(input string) error {
		if len(input) < 3 {
			return errors.New("username must have more than 3 characters")
		}
		return nil
	}
	return PromptForString("Username ", validate)
}

func PromptForLabels() ([]string, error) {
	validate := func(input string) error {
		return nil
	}
	lInput, err := PromptForString("Labels", validate)
	if err != nil {
		return nil, nil
	}
	l := strings.Split(lInput, ",")
	if len(l) == 0 {
		return nil, nil
	}
	return l, nil
}


func PromptForNewMPassword() (string, error) {
	validate := func(input string) error {
		if len(input) < 6 {
			return errors.New("new master password must have more than 6 characters")
		}
		return nil
	}
	return promptForPassword("New Master password ", validate)
}

func PromptForMPassword() (string, error) {
	validate := func(input string) error {
		if len(input) < 6 {
			return errors.New("master password must have more than 6 characters")
		}
		return nil
	}
	return promptForPassword("Master password ", validate)
}

// PromptForPassword function prompt for password and returns the input
func PromptForPassword() (string, error) {
	validate := func(input string) error {
		if len(input) < 6 {
			return errors.New("password must have more than 6 characters")
		}
		return nil
	}
	return promptForPassword("Password ", validate)
}
