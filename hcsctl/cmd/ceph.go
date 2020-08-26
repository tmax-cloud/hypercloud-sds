package cmd

import (
	"github.com/golang/glog"

	"hypercloud-sds/hcsctl/pkg/rook"
	"os"

	"github.com/spf13/cobra"
)

var cephCmd = &cobra.Command{
	Use:       "ceph",
	Short:     "ceph 의 상태를 확인하거나 접근합니다.",
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{statusCmd.Use, execCmd.Use},
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "ceph 의 현재 상태를 조회합니다.",
	PreRun: func(cmd *cobra.Command, args []string) {
		if isAvailable, err := rook.IsAvailableToolbox(); err != nil || !isAvailable {
			glog.Error("There isn't any available rook-ceph-toolbox pod in current k8s cluster.")
			glog.Error("Please check the rook-ceph-toolbox pod first.")
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := rook.Status()
		if err != nil {
			panic(err)
		}
	},
}

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "arguments 로 주어지는 ceph 명령을 수행합니다.",
	Args:  cobra.MinimumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		if isAvailable, err := rook.IsAvailableToolbox(); err != nil || !isAvailable {
			glog.Error("There isn't any available rook-ceph-toolbox pod in current k8s cluster.")
			glog.Error("Please check the rook-ceph-toolbox pod first.")
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := rook.Exec(args)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(cephCmd)
	cephCmd.AddCommand(statusCmd)
	cephCmd.AddCommand(execCmd)
}
