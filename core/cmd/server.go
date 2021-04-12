package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"onesite/common/config"
	"onesite/common/log"
	"onesite/core/dao"
	"onesite/core/server"
)

var (
	serverCfgFile string
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run server",
	Long:  "Run http server locally.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if serverCfgFile != "" {
			err := config.Init(serverCfgFile)
			if err != nil {
				return err
			}
			fmt.Printf("Run server with config file: %s\n", serverCfgFile)

			err = log.InitLogger()
			if err != nil {
				return err
			}

			err = dao.InitDao()
			if err != nil {
				return err
			}

			err = server.RunServer()
			if err != nil {
				return err
			}
		} else {
			_ = cmd.Help()
			os.Exit(0)
		}
		return nil
	},
}

func init() {
	serverCmd.Flags().StringVarP(&serverCfgFile, "conf", "c", "", "config [toml]")
	_ = serverCmd.MarkFlagRequired("conf")

	rootCmd.AddCommand(serverCmd)
}
