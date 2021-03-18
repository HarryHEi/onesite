package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:   "onesite",
	Short: "run onesite",
	Long:  "Run onesite.",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
		os.Exit(0)
	},
}

func Execute() error {
	return rootCmd.Execute()
}
