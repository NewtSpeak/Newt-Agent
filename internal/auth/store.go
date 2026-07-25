package auth

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/OwlSpeak/Owl-Agent/internal/config"
)

// 配置文件中 refresh 占位：真实 token 在 OS keyring。
const refreshKeyringMarker = "keyring"

// TokenSet 内存中的 access + 持久化 refresh。
type TokenSet struct {
	AccessToken     string
	RefreshToken    string
	AccessExpiresAt time.Time
	Scope           string
	ServerURL       string
}

var (
	mu      sync.Mutex
	current *TokenSet
)

func SetCurrent(t *TokenSet) {
	mu.Lock()
	defer mu.Unlock()
	current = t
}

func Current() *TokenSet {
	mu.Lock()
	defer mu.Unlock()
	if current == nil {
		return nil
	}
	cp := *current
	return &cp
}

// LoadFromConfig 从 keyring（优先）或配置文件加载 refresh。
func LoadFromConfig() (*TokenSet, error) {
	f, err := config.Load()
	if err != nil {
		return nil, err
	}
	p := f.Active()
	if p.ServerURL == "" {
		return nil, nil
	}
	refresh := resolveRefresh(f.ActiveProfile, p.RefreshToken)
	if refresh == "" {
		return nil, nil
	}
	t := &TokenSet{
		RefreshToken: refresh,
		Scope:        p.Scope,
		ServerURL:    p.ServerURL,
	}
	SetCurrent(t)
	return t, nil
}

func resolveRefresh(profile, fileValue string) string {
	// 1) keyring
	if tok, err := loadRefreshKeyring(profile); err == nil && tok != "" {
		return tok
	}
	// 2) 明文文件（兼容旧版 / keyring 不可用）
	if fileValue != "" && fileValue != refreshKeyringMarker {
		return fileValue
	}
	return ""
}

func SaveSession(serverURL, refresh, scope, clientID string) error {
	f, err := config.Load()
	if err != nil {
		return err
	}
	name := f.ActiveProfile
	if name == "" {
		name = "default"
	}
	p := f.Active()
	p.ServerURL = config.NormalizeServerURL(serverURL)
	p.Scope = scope
	p.ClientID = clientID

	// 优先 keyring；失败则落盘（仍 0600）
	if err := storeRefreshKeyring(name, refresh); err == nil {
		p.RefreshToken = refreshKeyringMarker
	} else {
		p.RefreshToken = refresh
	}

	f.SetActive(name, p)
	if err := config.Save(f); err != nil {
		return err
	}
	SetCurrent(&TokenSet{
		RefreshToken: refresh,
		Scope:        scope,
		ServerURL:    p.ServerURL,
	})
	return nil
}

func ClearSession() error {
	f, err := config.Load()
	if err != nil {
		return err
	}
	name := f.ActiveProfile
	deleteRefreshKeyring(name)
	p := f.Active()
	p.RefreshToken = ""
	p.Scope = ""
	f.SetActive(name, p)
	SetCurrent(nil)
	return config.Save(f)
}

// EnsureAccess 若 access 将过期则用 refresh 换新。
func EnsureAccess(refreshFn func(server, refresh string) (access, newRefresh string, exp time.Time, scope string, err error)) (string, error) {
	t := Current()
	if t == nil {
		loaded, err := LoadFromConfig()
		if err != nil {
			return "", err
		}
		t = loaded
	}
	if t == nil || t.RefreshToken == "" {
		return "", fmt.Errorf("未登录：请先运行 owl login --server <url>")
	}
	if t.AccessToken != "" && time.Now().Before(t.AccessExpiresAt.Add(-30*time.Second)) {
		return t.AccessToken, nil
	}
	access, newRefresh, exp, scope, err := refreshFn(t.ServerURL, t.RefreshToken)
	if err != nil {
		return "", err
	}
	if newRefresh == "" {
		newRefresh = t.RefreshToken
	}
	if scope == "" {
		scope = t.Scope
	}
	_ = SaveSession(t.ServerURL, newRefresh, scope, "owl-cli")
	mu.Lock()
	if current != nil {
		current.AccessToken = access
		current.AccessExpiresAt = exp
		current.RefreshToken = newRefresh
		current.Scope = scope
	}
	mu.Unlock()
	return access, nil
}

// DebugPath 返回配置文件路径。
func DebugPath() string {
	dir, err := config.Dir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "config.json")
}

func WriteHint() string {
	home, _ := os.UserConfigDir()
	return strings.TrimSpace(home) + "/owl-agent/config.json"
}

// SessionMeta 供 status tool 展示。
func SessionMeta() map[string]any {
	f, err := config.Load()
	if err != nil {
		return map[string]any{"error": err.Error()}
	}
	p := f.Active()
	refreshPresent := resolveRefresh(f.ActiveProfile, p.RefreshToken) != ""
	return map[string]any{
		"profile":         f.ActiveProfile,
		"server_url":      p.ServerURL,
		"scope":           p.Scope,
		"client_id":       p.ClientID,
		"logged_in":       refreshPresent,
		"config_path":     DebugPath(),
		"token_storage":   StorageBackend(f.ActiveProfile, p.RefreshToken != "" && p.RefreshToken != refreshKeyringMarker),
		"keyring_service": keyringService,
	}
}
