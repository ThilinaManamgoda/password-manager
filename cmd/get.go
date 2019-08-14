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
	"github.com/password-manager/pkg/encrypt"
	"github.com/password-manager/pkg/passwords"
	"github.com/password-manager/pkg/utils"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
)

const ShowPassword  = "show-pass"

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a password",
	Long:  `Get a password`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := utils.Configuration()
		if err != nil {
			return err
		}
		encriptorFac := &encrypt.Factory{
			ID: config.EncryptorID,
		}
		mPassword, err := utils.GetFlagStringVal(cmd, MasterPassword)
		if err != nil {
			return err
		}
		passwordRepo := &passwords.PasswordRepository{
			MasterPassword: mPassword,
			Encryptor:      encriptorFac.GetEncryptor(),
			PasswordFile: utils.PasswordFile{
				File: config.PasswordFilePath,
			},
		}
		showPass, err := utils.GetFlagBoolVal(cmd, ShowPassword)
		if err != nil {
			return err
		}
		if ! isArgSValid(args) {
			return errors.New("Give password entry ID")
		}
		if ! isArgValid(args[0]) {
			return errors.New(fmt.Sprintf("Invalid ID: %s", args[0]))
		}
		err = passwordRepo.GetPassword(args[0], showPass)
		if err != nil {
			return err
		}
		return nil
	},
}

func isArgSValid(args []string) bool {
	return len(args) != 0
}

func isArgValid(arg string) bool {
	return arg != ""
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	getCmd.Flags().BoolP(ShowPassword, "s", false, "Print password to STDOUT")
}
