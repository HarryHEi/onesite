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
	migrateCfgFile string
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate DB",
	Long:  "Migrate database or other data storage",
	RunE: func(cmd *cobra.Command, args []string) error {
		if migrateCfgFile != "" {
			err := config.Init(migrateCfgFile)
			if err != nil {
				return err
			}
			fmt.Printf("Run migrate with config file: %s\n", migrateCfgFile)

			err = log.InitLogger()
			if err != nil {
				return err
			}

			err = dao.InitDao()
			if err != nil {
				return err
			}

			err = dao.Migrate()
			if err != nil {
				return err
			}

			fmt.Println("Migrate successfully.")
		} else {
			_ = cmd.Help()
			os.Exit(0)
		}
		return nil
	},
}

func init() {
	migrateCmd.Flags().StringVarP(&migrateCfgFile, "conf", "c", "", "config [toml]")
	_ = migrateCmd.MarkFlagRequired("conf")

	rootCmd.AddCommand(migrateCmd)
}
