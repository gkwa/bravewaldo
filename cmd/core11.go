package cmd

import (
	"github.com/gkwa/bravewaldo/core11"
	"github.com/spf13/cobra"
)

// core11Cmd represents the core11 command
var core11Cmd = &cobra.Command{
	Use:   "core11",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := core11.Main()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(core11Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// core11Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// core11Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
