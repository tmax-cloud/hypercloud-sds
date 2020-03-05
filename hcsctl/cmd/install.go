package cmd

import (
	"github.com/spf13/cobra"
	"hypercloud-storage/hcsctl/pkg/cdi"
	"hypercloud-storage/hcsctl/pkg/rook"
)

var installCmd = &cobra.Command{
	Use:     "install",
	Short:   "해당 인벤토리를 기반으로 hypercloud-storage를 설치합니다.",
	PreRunE: checkAndSetInventory,
	Run: func(cmd *cobra.Command, args [] string) {
		err := rook.Apply(inventoryPath)
		if err != nil {
			panic(err)
		}
		err = cdi.Apply(inventoryPath)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
