package cmd

import (
	core "github.com/gkwa/bravewaldo/core1"
	"github.com/spf13/cobra"
)

var core1Cmd = &cobra.Command{
	Use:   "core1",
	Short: "A brief description of your command",
	Long:  `A longer description that spans multiple lines and likely contains examples and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := LoggerFrom(cmd.Context())
		core.Example(logger)
	},
}

func init() {
	rootCmd.AddCommand(core1Cmd)
}
