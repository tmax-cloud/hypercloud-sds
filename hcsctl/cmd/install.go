package cmd

import (
	"hypercloud-storage/hcsctl/pkg/cdi"
	"hypercloud-storage/hcsctl/pkg/rook"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use: "install",
	Short: "해당 인벤토리를 기반으로 hypercloud-storage를 설치합니다.\n" +
		"\t\t   자세한 예시는 hcsctl create-inventory 명령을 수행하면 확인할 수 있습니다.\n",
	PreRunE: validateInventory,
	Run: func(cmd *cobra.Command, args []string) {
		err := rook.Apply(inventoryPath)
		if err != nil {
			panic(err)
		}

		if isCdiExist(inventoryPath) {
			err = cdi.Apply(inventoryPath)
			if err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
