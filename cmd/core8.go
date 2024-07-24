package cmd

import (
	"github.com/gkwa/bravewaldo/core8"
	"github.com/spf13/cobra"
)

// core8Cmd represents the core8 command
var core8Cmd = &cobra.Command{
	Use:   "core8",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		core8.Main()
	},
}

func init() {
	rootCmd.AddCommand(core8Cmd)
}
