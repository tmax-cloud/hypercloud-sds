package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"

	glogcobra "github.com/blocktop/go-glog-cobra"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hcsctl.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

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

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".hcsctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".hcsctl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
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
