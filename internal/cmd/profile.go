package cmd

import (
	"fmt"
	"sort"

	"github.com/OwlSpeak/Owl-Agent/internal/auth"
	"github.com/OwlSpeak/Owl-Agent/internal/config"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "多账号 / 多服务器 profile 管理",
	Long: `每个 profile 独立保存 server_url、scope 与 refresh token（keyring 分区）。

  owl profile list
  owl profile use work
  owl login --server https://a.example --profile work
  owl profile show
  owl profile delete old`,
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出全部 profile",
	Run: func(cmd *cobra.Command, args []string) {
		f, err := config.Load()
		if err != nil {
			fatal(err)
		}
		names := f.ListProfiles()
		sort.Strings(names)
		type row struct {
			Name      string `json:"name"`
			Active    bool   `json:"active"`
			ServerURL string `json:"server_url,omitempty"`
			Scope     string `json:"scope,omitempty"`
			LoggedIn  bool   `json:"logged_in"`
		}
		out := make([]row, 0, len(names))
		active := f.ActiveProfile
		for _, n := range names {
			p := f.Profiles[n]
			logged := p.RefreshToken != "" // 含 keyring 占位
			out = append(out, row{
				Name: n, Active: n == active, ServerURL: p.ServerURL,
				Scope: p.Scope, LoggedIn: logged,
			})
		}
		// 尚无任何 profile 时至少展示 default 占位
		if len(out) == 0 {
			out = append(out, row{Name: "default", Active: true})
		}
		printJSON(out)
	},
}

var profileUseCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "切换当前 profile（后续命令使用该会话）",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		f, err := config.Load()
		if err != nil {
			fatal(err)
		}
		f.UseProfile(name)
		if err := config.Save(f); err != nil {
			fatal(err)
		}
		// 清内存 token，强制按新 profile 重载
		auth.SetCurrent(nil)
		fmt.Printf("已切换到 profile %q\n", f.ActiveProfile)
		p := f.Active()
		if p.ServerURL != "" {
			fmt.Println("server:", p.ServerURL)
		} else {
			fmt.Println("提示: 该 profile 尚未 login，请运行 owl login --server ...")
		}
	},
}

var profileShowCmd = &cobra.Command{
	Use:   "show",
	Short: "显示当前 profile 与会话元数据",
	Run: func(cmd *cobra.Command, args []string) {
		printJSON(auth.SessionMeta())
	},
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "删除 profile（并尝试清除 keyring）",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		f, err := config.Load()
		if err != nil {
			fatal(err)
		}
		// 先切到目标清会话
		prev := f.ActiveProfile
		f.ActiveProfile = name
		_ = auth.ClearSession() // 清 keyring + 当前文件项
		// ClearSession 会 reload 并清 active profile 的 token；再删 profile
		f2, err := config.Load()
		if err != nil {
			fatal(err)
		}
		if err := f2.DeleteProfile(name); err != nil {
			// 恢复
			f2.ActiveProfile = prev
			_ = config.Save(f2)
			fatal(err)
		}
		if f2.ActiveProfile == name {
			f2.ActiveProfile = "default"
		}
		if err := config.Save(f2); err != nil {
			fatal(err)
		}
		auth.SetCurrent(nil)
		fmt.Printf("已删除 profile %q，当前为 %q\n", name, f2.ActiveProfile)
	},
}

func init() {
	profileCmd.AddCommand(profileListCmd, profileUseCmd, profileShowCmd, profileDeleteCmd)
}
