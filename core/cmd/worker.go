package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"onesite/common/config"
	"onesite/common/log"
	"onesite/core/dao"
	"onesite/core/worker"
)

var (
	workerCfgFile string
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Run worker",
	Long:  "Run worker locally.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if workerCfgFile != "" {
			err := config.Init(workerCfgFile)
			if err != nil {
				return err
			}
			fmt.Printf("Run worker with config file: %s\n", workerCfgFile)

			err = log.InitLogger()
			if err != nil {
				return err
			}

			err = dao.InitDao()
			if err != nil {
				return err
			}

			err = worker.RunWorker()
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
	workerCmd.Flags().StringVarP(&workerCfgFile, "conf", "c", "", "config [toml]")
	_ = workerCmd.MarkFlagRequired("conf")

	rootCmd.AddCommand(workerCmd)
}
