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
	"fmt"
	"github.com/ThilinaManamgoda/password-manager/pkg/inputs"
	"github.com/ThilinaManamgoda/password-manager/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const LengthFlag = "length"

// generatePasswordCmd represents the generatePassword command
var generatePasswordCmd = &cobra.Command{
	Use:   "generate-password",
	Short: "Generate a secure password",
	Long:  `Generate a secure password of 12(default) characters in length`,
	RunE: func(cmd *cobra.Command, args []string) error {
		l, err := inputs.GetFlagIntVal(cmd, LengthFlag)
		if err != nil {
			return errors.Wrapf(err, inputs.ErrMSGCannotGetFlag, LengthFlag)
		}
		password, err := utils.GeneratePassword(l)
		if err != nil {
			return err
		}
		fmt.Println(password)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generatePasswordCmd)
	generatePasswordCmd.Flags().IntP(LengthFlag, "l", 12, "Length of the password")
}
