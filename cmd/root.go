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
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

// rootCmd represents the base command when called without any subcommands
var kubeconfig string

func init() {

	rootCmd = &cobra.Command{
		Use:   "crt",
		Short: "CodeReady Toolchain",
	}
	home := homeDir()
	rootCmd.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.SetHelpCommand(helpCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(createResourcesCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
