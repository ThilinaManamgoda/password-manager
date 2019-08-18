// Copyright © 2019 Thilina Manamgoda
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this fileio except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


// Package input handles the user interactions
package inputs

import (
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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

// PromptForPassword function prompt for password and returns the input
func PromptForPassword(label string, validate promptui.ValidateFunc) (string, error) {
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