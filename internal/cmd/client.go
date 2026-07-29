package cmd

import (
	"fmt"
	"time"

	"github.com/NewtSpeak/Newt-Agent/internal/api"
	"github.com/NewtSpeak/Newt-Agent/internal/auth"
	"github.com/NewtSpeak/Newt-Agent/internal/config"
	"github.com/NewtSpeak/Newt-Agent/internal/tools"
)

const defaultClientID = "owl-cli"

// ensureClient 返回带有效 access 的 API 客户端。
func ensureClient() *api.Client {
	c, err := makeClient()
	if err != nil {
		fatal(err)
	}
	return c
}

func makeClient() (*api.Client, error) {
	f, err := config.Load()
	if err != nil {
		return nil, err
	}
	server := f.Active().ServerURL
	if server == "" {
		return nil, fmt.Errorf("未配置服务器，请先 owl login --server <url>")
	}
	token, err := auth.EnsureAccess(func(serverURL, refresh string) (string, string, time.Time, string, error) {
		cl := api.New(serverURL)
		tr, err := cl.Refresh(refresh, defaultClientID)
		if err != nil {
			return "", "", time.Time{}, "", err
		}
		exp := tr.AccessExpiresAt
		if exp.IsZero() && tr.ExpiresIn > 0 {
			exp = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
		}
		return tr.AccessToken, tr.RefreshToken, exp, tr.Scope, nil
	})
	if err != nil {
		return nil, err
	}
	return api.New(server).WithToken(token), nil
}

func toolRegistry() *tools.Registry {
	return tools.NewRegistry(makeClient)
}

func runTool(name string, args map[string]any, yes bool) {
	reg := toolRegistry()
	result, err := reg.Call(name, args, yes)
	if err != nil {
		fatal(err)
	}
	text, err := tools.FormatResult(result)
	if err != nil {
		fatal(err)
	}
	fmt.Println(text)
}
