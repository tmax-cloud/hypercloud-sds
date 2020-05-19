package cmd

import (
	"errors"
	"flag"
	"fmt"
	"hypercloud-storage/hcsctl/pkg/kubectl"
	"os"
	"path/filepath"

	glogcobra "github.com/blocktop/go-glog-cobra"
	"github.com/spf13/cobra"
)

var kubeConfig string
var inventoryPath string

var rootCmd = &cobra.Command{
	Use:   "hcsctl",
	Short: "hypercloud-storage를 설치, 관리하기 위한 툴을 제공합니다.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	flag.Parse()
	cobra.OnInitialize(initConfig)

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().StringVar(&kubeConfig, "kubeconfig", filepath.Join(home, ".kube", "config"),
		"(optional) the location of kubeConfig file, default is $HOME/.kube/config")

	kubectl.KubeConfig = &kubeConfig

	glogcobra.Init(rootCmd)
}

// initConfig reads in config file and ENV variables if set.
// TODO: Should return error (ErrorHandling)
func initConfig() {
	err := glogcobra.Parse(rootCmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = flag.Set("logtostderr", "true")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func checkAndSetInventory(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("inventoryPath가 주어지지 않았습니다")
	}

	inventoryPath = args[0]

	return nil
}

func checkInventoryName(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("inventoryName 이 주어지지 않았습니다. 다음 형식으로 입력하세요 : " +
			"\n" + " - hcsctl create-inventory {$inventoryName}")
	}

	return nil
}
