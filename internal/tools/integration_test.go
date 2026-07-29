package tools_test

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/NewtSpeak/Newt-Agent/internal/api"
	"github.com/NewtSpeak/Newt-Agent/internal/tools"
	"github.com/gin-gonic/gin"
)

// TestToolsAgainstMockServer 不依赖真实 Newt-Server：用本地 mock 验证工具 HTTP 路径与解析。
func TestToolsAgainstMockServer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/oauth/v1/userinfo", func(c *gin.Context) {
		c.JSON(200, gin.H{"sub": "u1", "username": "alice", "system_admin": false})
	})
	r.GET("/gapi/v1/users/@me/guilds", func(c *gin.Context) {
		c.JSON(200, []gin.H{{"id": "g1", "name": "Test Guild"}})
	})
	r.GET("/gapi/v1/guilds/:gid/channels", func(c *gin.Context) {
		c.JSON(200, []gin.H{{"id": "c1", "name": "general", "type": "TEXT"}})
	})
	r.GET("/gapi/v1/search/messages", func(c *gin.Context) {
		if c.Query("q") == "" {
			c.JSON(400, gin.H{"error": gin.H{"code": "INVALID_QUERY"}})
			return
		}
		c.JSON(200, gin.H{"messages": []any{}, "total": 0})
	})
	r.GET("/gapi/v1/users/@me/relationships", func(c *gin.Context) {
		c.JSON(200, gin.H{"relationships": []any{}})
	})
	r.GET("/gapi/v1/users/@me/sticker-packs", func(c *gin.Context) {
		c.JSON(200, gin.H{"packs": []any{}})
	})
	srv := httptest.NewServer(r)
	defer srv.Close()

	reg := tools.NewRegistry(func() (*api.Client, error) {
		return api.New(srv.URL).WithToken("test-token"), nil
	})

	cases := []struct {
		name string
		args map[string]any
	}{
		{"whoami", nil},
		{"guilds.list", nil},
		{"channels.list", map[string]any{"guild_id": "g1"}},
		{"messages.search", map[string]any{"q": "hello", "guild_id": "g1", "limit": 10}},
		{"social.friends.list", nil},
		{"stickers.packs.list", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := reg.Call(tc.name, tc.args, false)
			if err != nil {
				t.Fatalf("%s: %v", tc.name, err)
			}
			if res == nil {
				t.Fatal("nil result")
			}
		})
	}
}

// TestOAuthDeviceFlowE2E 在设置 TEST_OWL_SERVER 时对真实服务器跑 device flow 冒烟（可选）。
// 需要可登录的测试账号不在此自动化（device 需人工同意），故仅检测 device/code 端点存活。
func TestOAuthDeviceCodeEndpointLive(t *testing.T) {
	base := os.Getenv("TEST_OWL_SERVER")
	if base == "" {
		t.Skip("TEST_OWL_SERVER not set")
	}
	client := api.New(base)
	dc, err := client.RequestDeviceCode("owl-cli", "profile gapi.full offline_access")
	if err != nil {
		t.Fatalf("device/code: %v", err)
	}
	if dc.DeviceCode == "" || dc.UserCode == "" {
		t.Fatalf("empty codes: %+v", dc)
	}
	// 立刻 poll 应 pending
	_, err = client.PollToken("owl-cli", dc.DeviceCode)
	if err == nil {
		t.Fatal("expected pending error")
	}
	if ae, ok := err.(*api.APIError); ok {
		if ae.Code != "authorization_pending" && ae.Code != "slow_down" {
			t.Logf("poll err code=%s msg=%s", ae.Code, ae.Message)
		}
	} else {
		t.Logf("poll err: %v", err)
	}
}
