package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"onesite/core/dao"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate DB",
	Long:  "Migrate database or other data storage",
	RunE: func(cmd *cobra.Command, args []string) error {
		dao, err := dao.NewDao()
		if err != nil {
			return err
		}

		err = dao.Migrate()
		if err != nil {
			return err
		}

		fmt.Println("Migrate successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
