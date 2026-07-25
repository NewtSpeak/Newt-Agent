package tools

import (
	"testing"

	"github.com/OwlSpeak/Owl-Agent/internal/api"
)

func TestRegistryListsCoreTools(t *testing.T) {
	r := NewRegistry(func() (*api.Client, error) {
		return api.New("http://example.invalid"), nil
	})
	list := r.List()
	if len(list) < 55 {
		t.Fatalf("expected many tools, got %d", len(list))
	}
	need := []string{
		"whoami", "guilds.list", "channels.list", "channels.create",
		"channels.delete", "roles.list", "members.kick", "messages.send",
		"invites.create", "restrictions.list", "audit.list", "voice.disconnect",
		"platform.users.list", "platform.sfu.nodes",
		"stickers.packs.list", "social.friends.list", "social.dm.list",
	}
	set := map[string]bool{}
	for _, d := range list {
		set[d.Name] = true
		if d.InputSchema == nil {
			t.Fatalf("tool %s missing inputSchema", d.Name)
		}
	}
	for _, n := range need {
		if !set[n] {
			t.Fatalf("missing tool %s", n)
		}
	}
}

func TestDestructiveRequiresConfirm(t *testing.T) {
	r := NewRegistry(func() (*api.Client, error) {
		return api.New("http://example.invalid"), nil
	})
	_, err := r.Call("channels.delete", map[string]any{"channel_id": "x"}, false)
	if err == nil {
		t.Fatal("expected confirm error")
	}
}
