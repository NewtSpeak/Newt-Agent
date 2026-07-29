package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/NewtSpeak/Newt-Agent/internal/api"
	"github.com/NewtSpeak/Newt-Agent/internal/auth"
	"github.com/NewtSpeak/Newt-Agent/internal/config"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "检查登录、网络与 OAuth 端点健康状况",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("owl doctor")
		fmt.Println("----------")

		meta := auth.SessionMeta()
		fmt.Printf("config:   %v\n", meta["config_path"])
		fmt.Printf("profile:  %v\n", meta["profile"])
		fmt.Printf("storage:  %v\n", meta["token_storage"])
		fmt.Printf("logged_in:%v\n", meta["logged_in"])

		server, _ := meta["server_url"].(string)
		if server == "" {
			f, _ := config.Load()
			server = f.Active().ServerURL
		}
		if server == "" {
			fmt.Println("server:   (未配置) 请 owl login --server <url>")
			os.Exit(1)
		}
		fmt.Printf("server:   %s\n", server)

		client := &http.Client{Timeout: 8 * time.Second}
		// healthz
		resp, err := client.Get(server + "/healthz")
		if err != nil {
			fmt.Printf("healthz:  FAIL %v\n", err)
		} else {
			_ = resp.Body.Close()
			fmt.Printf("healthz:  %s\n", resp.Status)
		}

		// OAuth discovery
		resp, err = client.Get(server + "/oauth/v1/.well-known/oauth-authorization-server")
		if err != nil {
			fmt.Printf("oauth:    FAIL %v\n", err)
		} else {
			_ = resp.Body.Close()
			fmt.Printf("oauth:    discovery %s\n", resp.Status)
		}

		if meta["logged_in"] == true {
			cl, err := makeClient()
			if err != nil {
				fmt.Printf("whoami:   FAIL %v\n", err)
				os.Exit(1)
			}
			info, err := cl.UserInfo()
			if err != nil {
				fmt.Printf("whoami:   FAIL %v\n", err)
				// 尝试 refresh 路径已在 makeClient
				os.Exit(1)
			}
			fmt.Printf("whoami:   %s (admin=%v)\n", info.Username, info.SystemAdmin)
			fmt.Printf("scope:    %s\n", info.Scope)
		} else {
			fmt.Println("whoami:   skipped (未登录)")
		}

		// 轻量 gapi 探测
		if meta["logged_in"] == true {
			cl := api.New(server)
			// 用 ensure 的 token
			c2, err := makeClient()
			if err == nil {
				cl = c2
			}
			raw, err := cl.Gapi(http.MethodGet, "/users/@me/guilds", nil, nil)
			if err != nil {
				fmt.Printf("gapi:     FAIL %v\n", err)
			} else {
				fmt.Printf("gapi:     guilds OK (%d bytes)\n", len(raw))
			}
		}

		fmt.Println("----------")
		fmt.Println("doctor 完成。")
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
