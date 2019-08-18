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
	"github.com/spf13/cobra/doc"
	"os"
	"path"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// tool version. Should be initialized at build time
var Version string
var IsGenerateDoc string

// InteractiveMode flag
const InteractiveMode = "interactive"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "password-manager",
	Short: "A local Password Manager",
	Long:  `A local password manager`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if IsGenerateDoc == "true" {
		wd, err:= os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = doc.GenMarkdownTree(rootCmd,path.Join(wd))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "c", "config file (default is $HOME/.password-manager.yaml)")
	rootCmd.PersistentFlags().StringP(MasterPassword, "m", "", "Master password")
	rootCmd.PersistentFlags().BoolP(InteractiveMode, "i", false, "Enable interactive mode")
	//addCmd.Flags().StringP(mPassword, "m", "", "Master password")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Version = Version

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".password-manager" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".password-manager")
	}
	viper.SetEnvPrefix("PASSWORD_MANAGER")
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
