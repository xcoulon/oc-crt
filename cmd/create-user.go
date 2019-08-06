package cmd

import (
	"github.com/spf13/cobra"
)

var createUserCmd = &cobra.Command{
	Use:               "create-user",
	Short:             "creates a new CodeReady Toolchain user on the host cluster",
	PersistentPreRun:  func(cmd *cobra.Command, args []string) {},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {},
	Args:              cobra.ExactArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		// if err := o.Complete(c, args); err != nil {
		// 	return err
		// }
		// if err := o.Validate(); err != nil {
		// 	return err
		// }
		// if err := o.Run(); err != nil {
		// 	return err
		// }

		return nil
	},
}
