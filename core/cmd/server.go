package cmd

import (
	"github.com/spf13/cobra"

	"onesite/core/server"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run server",
	Long:  "Run http server locally.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return server.RunServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
