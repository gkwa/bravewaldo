/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/gkwa/bravewaldo/core10"
	"github.com/spf13/cobra"
)

// core10Cmd represents the core10 command
var core10Cmd = &cobra.Command{
	Use:   "core10",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := LoggerFrom(cmd.Context())
		core10.Main(logger)
	},
}

func init() {
	rootCmd.AddCommand(core10Cmd)
}
