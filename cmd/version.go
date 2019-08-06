package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// BinaryName name of the binary (set by Makefile)
	BinaryName = "oc-crt"
	// BuildCommit lastest build commit (set by Makefile)
	BuildCommit = "git+master"
	// BuildTag if the `BuildCommit` matches a tag
	BuildTag = "+git"
)

var versionCmd = &cobra.Command{
	Use:               "version",
	Short:             "print the version of this plugin",
	PersistentPreRun:  func(cmd *cobra.Command, args []string) {},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {},
	Run: func(c *cobra.Command, args []string) {
		if BuildTag != "" {
			fmt.Printf("%s version %s\n", BinaryName, BuildTag)
		} else {
			fmt.Printf("%s version %s\n", BinaryName, BuildCommit)
		}
	},
}
