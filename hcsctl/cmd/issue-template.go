package cmd

import (
	"hypercloud-sds/hcsctl/pkg/rook"

	"github.com/spf13/cobra"
)

var issueTemplateCmd = &cobra.Command{
	Use:   "issue-template",
	Short: "issue-template을 생성합니다.",
	Run: func(cmd *cobra.Command, args []string) {
		err := rook.GetIssueTemplate()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(issueTemplateCmd)
}
