package cmd

import (
	"github.com/spf13/cobra"
)

var guildsCmd = &cobra.Command{
	Use:   "guilds",
	Short: "服务器（Guild）相关命令",
}

var guildsListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出当前用户加入的服务器",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("guilds.list", map[string]any{}, flagYes)
	},
}

func init() {
	guildsCmd.AddCommand(guildsListCmd)
}
