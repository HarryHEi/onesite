package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"onesite/core/dao"
)

var (
	initUsername string
	initPassword string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init",
	Long:  "Create superuser if not exists.",
	RunE: func(cmd *cobra.Command, args []string) error {
		dao, err := dao.NewDao()
		if err != nil {
			return err
		}

		err = dao.CreateSuperuserIfNotExists(initUsername, initPassword)
		if err != nil {
			return err
		}
		fmt.Println("Super user created.")

		fmt.Println("Init successfully.")
		return nil
	},
}

func init() {
	initCmd.Flags().StringVarP(&initUsername, "username", "u", "admin", "username")
	initCmd.Flags().StringVarP(&initPassword, "password", "p", "admin", "password")
	_ = initCmd.MarkFlagRequired("conf")

	rootCmd.AddCommand(initCmd)
}
