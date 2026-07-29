package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "newt",
	Short: "NewtSpeak CLI — OAuth 用户委托 + AI skill / MCP 工具",
	Long: `newt 是 NewtSpeak 官方 Agent CLI。

登录后可管理服务器、频道、角色、成员与消息，并作为 AI skill / MCP 工具入口。
密码不会进入本工具：使用设备码或 PKCE 在 Desktop / Web 授权页完成登录。

  newt login --server https://...
  newt login --method pkce --client-origin https://web...
  newt doctor
  newt channels list --guild <id>
  newt mcp serve

Shell 补全：
  newt completion bash > /etc/bash_completion.d/owl
  newt completion powershell | Out-String | Invoke-Expression`,
	Version: Version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	bindYes(rootCmd)
	initAdminCommands()
	initExtraCommands()
	initSocialCommands()

	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(whoamiCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(profileCmd)
	rootCmd.AddCommand(guildsCmd)
	rootCmd.AddCommand(channelsCmd)
	rootCmd.AddCommand(rolesCmd)
	rootCmd.AddCommand(membersCmd)
	rootCmd.AddCommand(invitesCmd)
	rootCmd.AddCommand(messagesCmd)
	rootCmd.AddCommand(restrictionsCmd)
	rootCmd.AddCommand(auditCmd)
	rootCmd.AddCommand(voiceCmd)
	rootCmd.AddCommand(platformCmd)
	rootCmd.AddCommand(stickersCmd)
	rootCmd.AddCommand(socialCmd)
	rootCmd.AddCommand(toolsCmd)
	rootCmd.AddCommand(mcpCmd)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "错误:", err)
	os.Exit(1)
}
