package mcp

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/NewtSpeak/Newt-Agent/internal/api"
	"github.com/NewtSpeak/Newt-Agent/internal/tools"
)

func TestMCPInitializeResourcesPrompts(t *testing.T) {
	reg := tools.NewRegistry(func() (*api.Client, error) {
		return api.New("http://example.invalid"), nil
	})
	var out bytes.Buffer
	in := strings.NewReader(strings.Join([]string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`,
		`{"jsonrpc":"2.0","id":2,"method":"resources/list","params":{}}`,
		`{"jsonrpc":"2.0","id":3,"method":"prompts/list","params":{}}`,
		`{"jsonrpc":"2.0","id":4,"method":"prompts/get","params":{"name":"safe-ops"}}`,
		`{"jsonrpc":"2.0","id":5,"method":"resources/read","params":{"uri":"newtspeak://status"}}`,
	}, "\n") + "\n")

	s := &Server{Reg: reg, Reader: in, Writer: &out, ErrLog: &bytes.Buffer{}}
	if err := s.Serve(); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) < 5 {
		t.Fatalf("expected 5 responses, got %d: %s", len(lines), out.String())
	}
	var init map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &init); err != nil {
		t.Fatal(err)
	}
	result := init["result"].(map[string]any)
	caps := result["capabilities"].(map[string]any)
	if caps["resources"] == nil || caps["prompts"] == nil {
		t.Fatalf("missing resources/prompts caps: %#v", caps)
	}
	var promptsGet map[string]any
	if err := json.Unmarshal([]byte(lines[3]), &promptsGet); err != nil {
		t.Fatal(err)
	}
	if promptsGet["error"] != nil {
		t.Fatalf("prompts/get error: %v", promptsGet["error"])
	}
}
