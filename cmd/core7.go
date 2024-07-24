package cmd

import (
	"github.com/gkwa/bravewaldo/core7"
	"github.com/spf13/cobra"
)

// core7Cmd represents the core7 command
var core7Cmd = &cobra.Command{
	Use:   "core7",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		core7.Main()
	},
}

func init() {
	rootCmd.AddCommand(core7Cmd)
}
