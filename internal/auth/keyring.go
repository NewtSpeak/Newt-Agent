package auth

import (
	"fmt"
	"os"
	"strings"

	"github.com/zalando/go-keyring"
)

const keyringService = "newt-agent"

// keyring 不可用或失败时回退到配置文件（仍 0600）。
var keyringDisabled = strings.EqualFold(os.Getenv("NEWT_AGENT_NO_KEYRING"), "1") ||
	strings.EqualFold(os.Getenv("NEWT_AGENT_NO_KEYRING"), "true")

func keyringAccount(profile string) string {
	if profile == "" {
		profile = "default"
	}
	return "refresh:" + profile
}

func storeRefreshKeyring(profile, token string) error {
	if keyringDisabled || token == "" {
		return errKeyringSkip
	}
	if err := goKeyringSet(keyringService, keyringAccount(profile), token); err != nil {
		return err
	}
	return nil
}

func loadRefreshKeyring(profile string) (string, error) {
	if keyringDisabled {
		return "", errKeyringSkip
	}
	return goKeyringGet(keyringService, keyringAccount(profile))
}

func deleteRefreshKeyring(profile string) {
	if keyringDisabled {
		return
	}
	_ = goKeyringDelete(keyringService, keyringAccount(profile))
}

var errKeyringSkip = fmt.Errorf("keyring skipped")

// 可替换的函数变量，便于测试。
var (
	goKeyringSet    = keyring.Set
	goKeyringGet    = keyring.Get
	goKeyringDelete = keyring.Delete
)

// StorageBackend 描述 refresh 实际存放位置。
func StorageBackend(profile string, fileHasToken bool) string {
	if keyringDisabled {
		if fileHasToken {
			return "config-file"
		}
		return "none"
	}
	if _, err := loadRefreshKeyring(profile); err == nil {
		return "os-keyring"
	}
	if fileHasToken {
		return "config-file"
	}
	return "none"
}
