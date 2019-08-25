// +build doc

package main

import (
	"fmt"
	"github.com/spf13/cobra/doc"
	"github.com/ThilinaManamgoda/password-manager/cmd"
	"os"
	"path"
)

func main() {
	wd, err:= os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd := cmd.GetRootCMD()
	err = doc.GenMarkdownTree(rootCmd,path.Join(wd))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}