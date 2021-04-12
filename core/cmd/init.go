package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"onesite/common/config"
	"onesite/common/log"
	"onesite/core/dao"
)

var (
	initCfgFile  string
	initUsername string
	initPassword string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init",
	Long:  "Create superuser if not exists.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if initCfgFile != "" {
			err := config.Init(initCfgFile)
			if err != nil {
				return err
			}
			fmt.Printf("Run server with config file: %s\n", initCfgFile)

			err = log.InitLogger()
			if err != nil {
				return err
			}

			err = dao.InitDao()
			if err != nil {
				return err
			}

			return dao.CreateSuperuserIfNotExists(initUsername, initPassword)
		} else {
			_ = cmd.Help()
			os.Exit(0)
		}
		return nil
	},
}

func init() {
	initCmd.Flags().StringVarP(&initCfgFile, "conf", "c", "", "config [toml]")
	initCmd.Flags().StringVarP(&initUsername, "username", "u", "admin", "username")
	initCmd.Flags().StringVarP(&initPassword, "password", "p", "admin", "password")
	_ = initCmd.MarkFlagRequired("conf")

	rootCmd.AddCommand(initCmd)
}
