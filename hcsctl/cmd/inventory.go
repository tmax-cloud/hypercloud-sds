package cmd

import (
	"hypercloud-sds/hcsctl/pkg/inventory"

	"github.com/spf13/cobra"
)

var inventoryCreateCmd = &cobra.Command{
	Use: "create-inventory",
	Short: "hcsctl install 시 사용할 수 있는 {$inventoryName} 디렉토리를 '현재 경로'에 생성합니다.\n" +
		"\t\t   해당 디렉토리는 ./{$inventoryName}/rook/ 아래 rook 관련 yaml 파일,\n" +
		"\t\t   그리고 ./{$inventoryName}/cdi/ 아래 cdi 관련 yaml 파일을 담고 있습니다.\n",
	PreRunE: checkInventoryName,
	Args:    cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inventoryName := args[0]
		includingCdiFlag, _ := cmd.Flags().GetBool("include-cdi")

		err := inventory.Create(inventoryName, includingCdiFlag)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	inventoryCreateCmd.Flags().Bool("include-cdi", false, "CDI Installation Feature. Bool type")
	rootCmd.AddCommand(inventoryCreateCmd)
}
