package cmd

import (
	"github.com/spf13/cobra"
	"hypercloud-storage/hcsctl/pkg/cdi"
	"hypercloud-storage/hcsctl/pkg/rook"
)

var uninstallCmd = &cobra.Command{
	Use:     "uninstall",
	Short:   "해당 인벤토리의 hypercloud-storage를 제거합니다.",
	PreRunE: checkAndSetInventory,
	Run: func(cmd *cobra.Command, args []string) {
		err := cdi.Delete(inventoryPath)
		if err != nil {
			panic(err)
		}
		err = rook.Delete(inventoryPath)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
