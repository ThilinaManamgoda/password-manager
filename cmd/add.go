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
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new password",
	Long:  `Add a new password`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("add called")
		encriptorFac := &encrypt.Factory{

		}
		fileManager := &utils.PasswordFile{
			File: "/Users/wso2/.go/src/github.com/password-manager/test1",
		}


		passwordRepo := &passwords.PasswordRepository{
			MasterPassword: "test",
			Encryptor:encriptorFac.GetEncryptor(),
			PasswordFile: utils.PasswordFile{
				File: "/Users/wso2/.go/src/github.com/password-manager/test1",
			},
		}

		err:=passwordRepo.Add("test","test", nil)
		if err != nil {
			panic(err)
		}


		encriptor := encriptorFac.GetEncryptor()
		//encryptedData, _ := encriptor.Encrypt([]byte("test"), "test")
		//fmt.Println(encryptedData)
		////fileManager.StorePasswords(encryptedData)
		encryptedData, _ := fileManager.GetPasswords()
		decryptedData, _ := encriptor.Decrypt(encryptedData, "test")
		fmt.Println(string(decryptedData))
		passwordRepo.GetPassword("test", true)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	addCmd.Flags().StringP("password", "p", "", "Master password")
}
