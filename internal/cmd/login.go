package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/NewtSpeak/Newt-Agent/internal/api"
	"github.com/NewtSpeak/Newt-Agent/internal/auth"
	"github.com/NewtSpeak/Newt-Agent/internal/config"
	"github.com/spf13/cobra"
)

var (
	loginServer   string
	loginPlatform bool
	loginNoOpen   bool
	loginProfile  string
	loginMethod   string // device | pkce
	loginClientOrigin string
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "OAuth 登录（device 设备码 或 pkce 浏览器回调）",
	Run: func(cmd *cobra.Command, args []string) {
		if loginProfile != "" {
			f, err := config.Load()
			if err != nil {
				fatal(err)
			}
			f.UseProfile(loginProfile)
			if err := config.Save(f); err != nil {
				fatal(err)
			}
			auth.SetCurrent(nil)
			fmt.Println("使用 profile:", loginProfile)
		}
		if loginServer == "" {
			f, _ := config.Load()
			if f.Active().ServerURL != "" {
				loginServer = f.Active().ServerURL
			}
		}
		if loginServer == "" {
			fatal(fmt.Errorf("请指定 --server https://你的-owl-server"))
		}
		server := config.NormalizeServerURL(loginServer)
		scope := "openid profile gapi.full offline_access"
		if loginPlatform {
			scope += " platform.read platform.admin"
		}
		method := loginMethod
		if method == "" {
			method = "device"
		}
		switch method {
		case "device":
			loginDevice(server, scope)
		case "pkce":
			loginPKCE(server, scope)
		default:
			fatal(fmt.Errorf("未知 --method %q（device|pkce）", method))
		}
	},
}

func loginDevice(server, scope string) {
	client := api.New(server)
	dc, err := client.RequestDeviceCode(defaultClientID, scope)
	if err != nil {
		fatal(err)
	}
	interval := dc.Interval
	if interval < 1 {
		interval = 5
	}

	fmt.Println()
	fmt.Println("请在 NewtSpeak 客户端或浏览器完成授权：")
	fmt.Println()
	fmt.Println("  设备码:", dc.UserCode)
	fmt.Println("  打开:  ", dc.VerificationURIComplete)
	fmt.Println()
	fmt.Println("提示：已登录的 Desktop 可通过深链 newtspeak://oauth/device 打开。")
	fmt.Println("勿将密码或 refresh token 发给 AI。")
	fmt.Println()

	if !loginNoOpen {
		_ = openBrowser(dc.VerificationURIComplete)
		deep := fmt.Sprintf("newtspeak://oauth/device?user_code=%s&server=%s", dc.UserCode, server)
		_ = openBrowser(deep)
	}

	deadline := time.Now().Add(time.Duration(dc.ExpiresIn) * time.Second)
	for time.Now().Before(deadline) {
		time.Sleep(time.Duration(interval) * time.Second)
		tr, err := client.PollToken(defaultClientID, dc.DeviceCode)
		if err != nil {
			if ae, ok := err.(*api.APIError); ok {
				switch ae.Code {
				case "authorization_pending":
					fmt.Fprint(os.Stderr, ".")
					continue
				case "slow_down":
					interval += 5
					fmt.Fprint(os.Stderr, "+")
					continue
				case "access_denied":
					fatal(fmt.Errorf("用户拒绝了授权"))
				case "expired_token":
					fatal(fmt.Errorf("设备码已过期，请重新 owl login"))
				}
			}
			fatal(err)
		}
		if tr == nil || tr.AccessToken == "" {
			continue
		}
		fmt.Fprintln(os.Stderr)
		finishLogin(server, tr)
		return
	}
	fatal(fmt.Errorf("等待授权超时"))
}

func loginPKCE(server, scope string) {
	verifier, challenge, err := api.GeneratePKCE()
	if err != nil {
		fatal(err)
	}
	state, err := api.RandomState()
	if err != nil {
		fatal(err)
	}
	redirectURI, stop, codeCh, err := api.StartLoopback(state)
	if err != nil {
		fatal(err)
	}
	defer stop()

	origin := loginClientOrigin
	if origin == "" {
		f, _ := config.Load()
		origin = f.Active().ClientOrigin
	}
	authURL := api.BuildAuthorizeURL(origin, server, defaultClientID, redirectURI, scope, challenge, state)

	fmt.Println()
	fmt.Println("PKCE 登录：请在浏览器 / Desktop 完成授权")
	fmt.Println()
	fmt.Println("  打开:", authURL)
	fmt.Println()
	if !loginNoOpen {
		_ = openBrowser(authURL)
	}

	select {
	case res := <-codeCh:
		if res.Err != nil {
			fatal(res.Err)
		}
		client := api.New(server)
		tr, err := client.ExchangeAuthCode(defaultClientID, res.Code, verifier, redirectURI)
		if err != nil {
			fatal(err)
		}
		finishLogin(server, tr)
	case <-time.After(5 * time.Minute):
		fatal(fmt.Errorf("等待授权超时"))
	}
}

func finishLogin(server string, tr *api.TokenResponse) {
	if err := auth.SaveSession(server, tr.RefreshToken, tr.Scope, defaultClientID); err != nil {
		fatal(err)
	}
	exp := tr.AccessExpiresAt
	if exp.IsZero() && tr.ExpiresIn > 0 {
		exp = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	}
	auth.SetCurrent(&auth.TokenSet{
		AccessToken:     tr.AccessToken,
		RefreshToken:    tr.RefreshToken,
		AccessExpiresAt: exp,
		Scope:           tr.Scope,
		ServerURL:       server,
	})
	// 记住 client origin（若指定）
	if loginClientOrigin != "" {
		f, err := config.Load()
		if err == nil {
			p := f.Active()
			p.ClientOrigin = config.NormalizeServerURL(loginClientOrigin)
			f.SetActive(f.ActiveProfile, p)
			_ = config.Save(f)
		}
	}
	fmt.Println("登录成功。")
	fmt.Println("scope:", tr.Scope)
	fmt.Println("配置:", auth.DebugPath())
}

func init() {
	loginCmd.Flags().StringVar(&loginServer, "server", "", "Newt-Server 基址")
	loginCmd.Flags().BoolVar(&loginPlatform, "platform", false, "申请平台 scope")
	loginCmd.Flags().BoolVar(&loginNoOpen, "no-open", false, "不自动打开浏览器")
	loginCmd.Flags().StringVar(&loginProfile, "profile", "", "写入指定 profile")
	loginCmd.Flags().StringVar(&loginMethod, "method", "device", "device | pkce")
	loginCmd.Flags().StringVar(&loginClientOrigin, "client-origin", "", "用户 Web/Desktop 前端基址（PKCE 授权页）")
}

func openBrowser(url string) error {
	if url == "" {
		return nil
	}
	var c *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		c = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		c = exec.Command("open", url)
	default:
		c = exec.Command("xdg-open", url)
	}
	return c.Start()
}
