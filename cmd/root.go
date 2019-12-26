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
	"fmt"
	"github.com/ThilinaManamgoda/password-manager/pkg/config"
	"github.com/ThilinaManamgoda/password-manager/pkg/inputs"
	"github.com/spf13/cobra"
	"os"
)

// Version of the password manager. Should be initialized at build time
var Version string

// InteractiveMode flag
const InteractiveMode = "interactive"

// rootCmd represents the base command when called without any sub-commands
var rootCmd = &cobra.Command{
	Use:   "password-manager",
	Short: "A local Password Manager",
	Long:  `A local password manager`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.Init)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringP(inputs.FlagMasterPassword, "m", "", "Master password")
	//addCmd.Flags().StringP(mPassword, "m", "", "Master password")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Version = Version

}

// GetRootCMD returns the Root CMD struct
func GetRootCMD() *cobra.Command {
	return rootCmd
}
