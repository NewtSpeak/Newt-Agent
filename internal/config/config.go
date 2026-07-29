package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Profile struct {
	ServerURL    string `json:"server_url"`
	ClientOrigin string `json:"client_origin,omitempty"`
	// RefreshToken 明文仅在无法使用 keyring 时落入此文件（权限应 0600）。
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
	ClientID     string `json:"client_id,omitempty"`
}

type File struct {
	ActiveProfile string             `json:"active_profile"`
	Profiles      map[string]Profile `json:"profiles"`
}

func Dir() (string, error) {
	home, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, "newt-agent")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return dir, nil
}

func path() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

func Load() (File, error) {
	p, err := path()
	if err != nil {
		return File{}, err
	}
	data, err := os.ReadFile(p)
	if errors.Is(err, os.ErrNotExist) {
		return File{ActiveProfile: "default", Profiles: map[string]Profile{}}, nil
	}
	if err != nil {
		return File{}, err
	}
	var f File
	if err := json.Unmarshal(data, &f); err != nil {
		return File{}, err
	}
	if f.Profiles == nil {
		f.Profiles = map[string]Profile{}
	}
	if f.ActiveProfile == "" {
		f.ActiveProfile = "default"
	}
	return f, nil
}

func Save(f File) error {
	p, err := path()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o600)
}

func (f *File) Active() Profile {
	if f.Profiles == nil {
		return Profile{}
	}
	return f.Profiles[f.ActiveProfile]
}

func (f *File) SetActive(name string, p Profile) {
	if f.Profiles == nil {
		f.Profiles = map[string]Profile{}
	}
	if name == "" {
		name = "default"
	}
	f.ActiveProfile = name
	f.Profiles[name] = p
}

// UseProfile 切换当前 profile 名称（不改内容）；不存在则创建空 profile。
func (f *File) UseProfile(name string) {
	if f.Profiles == nil {
		f.Profiles = map[string]Profile{}
	}
	name = strings.TrimSpace(name)
	if name == "" {
		name = "default"
	}
	if _, ok := f.Profiles[name]; !ok {
		f.Profiles[name] = Profile{}
	}
	f.ActiveProfile = name
}

// DeleteProfile 删除命名 profile；不可删除当前唯一/正在使用的需调用方处理。
func (f *File) DeleteProfile(name string) error {
	name = strings.TrimSpace(name)
	if name == "" || name == "default" && len(f.Profiles) <= 1 {
		return errors.New("不能删除唯一的 default profile")
	}
	if f.Profiles == nil {
		return errors.New("profile 不存在")
	}
	if _, ok := f.Profiles[name]; !ok {
		return errors.New("profile 不存在: " + name)
	}
	delete(f.Profiles, name)
	if f.ActiveProfile == name {
		f.ActiveProfile = "default"
		if _, ok := f.Profiles["default"]; !ok {
			// 任选一个
			for k := range f.Profiles {
				f.ActiveProfile = k
				break
			}
		}
	}
	return nil
}

// ListProfiles 返回名称列表（含是否 active 由调用方判断）。
func (f *File) ListProfiles() []string {
	out := make([]string, 0, len(f.Profiles))
	for k := range f.Profiles {
		out = append(out, k)
	}
	return out
}

func NormalizeServerURL(raw string) string {
	s := strings.TrimSpace(raw)
	s = strings.TrimRight(s, "/")
	return s
}

