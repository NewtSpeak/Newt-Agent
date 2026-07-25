package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/OwlSpeak/Owl-Agent/internal/api"
	"github.com/OwlSpeak/Owl-Agent/internal/auth"
	"github.com/OwlSpeak/Owl-Agent/internal/config"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "吊销本地 refresh 并清除会话",
	Run: func(cmd *cobra.Command, args []string) {
		f, err := config.Load()
		if err != nil {
			fatal(err)
		}
		p := f.Active()
		if p.ServerURL != "" && p.RefreshToken != "" {
			cl := api.New(p.ServerURL)
			_ = cl.Revoke(p.RefreshToken)
		}
		if err := auth.ClearSession(); err != nil {
			fatal(err)
		}
		fmt.Println("已退出登录。")
	},
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "显示当前 OAuth 用户信息",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("whoami", map[string]any{}, flagYes)
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "显示登录状态与配置路径",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("status", map[string]any{}, flagYes)
	},
}

func marshalIndent(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

func printJSON(v any) {
	b, err := marshalIndent(v)
	if err != nil {
		fatal(err)
	}
	fmt.Println(string(b))
}
