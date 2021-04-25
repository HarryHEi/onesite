package cmd

import (
	"github.com/spf13/cobra"

	"onesite/core/dao"
	"onesite/core/worker"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Run worker",
	Long:  "Run worker locally.",
	RunE: func(cmd *cobra.Command, args []string) error {
		dao, err := dao.NewDao()
		if err != nil {
			return err
		}

		w := worker.NewWorker(dao)
		w.Run()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
}
