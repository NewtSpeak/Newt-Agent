package tools

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/OwlSpeak/Owl-Agent/internal/api"
)

// ClientFactory 在每次调用时提供已认证的 API 客户端。
type ClientFactory func() (*api.Client, error)

// Def 工具定义（CLI / MCP / Skill 共用）。
type Def struct {
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Destructive  bool           `json:"destructive,omitempty"`
	InputSchema  map[string]any `json:"inputSchema"`
	CLIHint      string         `json:"cli_hint,omitempty"`
	run          func(c *api.Client, args map[string]any) (any, error)
}

// Registry 全部工具。
type Registry struct {
	factory ClientFactory
	byName  map[string]*Def
	order   []string
}

func NewRegistry(factory ClientFactory) *Registry {
	r := &Registry{factory: factory, byName: map[string]*Def{}}
	r.registerAll()
	return r
}

func (r *Registry) add(d *Def) {
	r.byName[d.Name] = d
	r.order = append(r.order, d.Name)
}

func (r *Registry) List() []Def {
	out := make([]Def, 0, len(r.order))
	for _, name := range r.order {
		d := r.byName[name]
		out = append(out, Def{
			Name: d.Name, Description: d.Description, Destructive: d.Destructive,
			InputSchema: d.InputSchema, CLIHint: d.CLIHint,
		})
	}
	return out
}

func (r *Registry) Get(name string) (*Def, bool) {
	d, ok := r.byName[name]
	return d, ok
}

// Call 执行工具；destructive 工具要求 args.confirm == true（除非 force）。
func (r *Registry) Call(name string, args map[string]any, force bool) (any, error) {
	d, ok := r.byName[name]
	if !ok {
		names := r.order
		sort.Strings(names)
		return nil, fmt.Errorf("未知工具 %q；可用: %s", name, strings.Join(names, ", "))
	}
	if args == nil {
		args = map[string]any{}
	}
	if d.Destructive && !force {
		if !truthy(args["confirm"]) {
			return nil, fmt.Errorf("危险操作 %s 需要 confirm=true（或 CLI --yes）", name)
		}
	}
	// status 可在未登录时查看本地配置；其余工具需要有效会话。
	var c *api.Client
	if name != "status" {
		var err error
		c, err = r.factory()
		if err != nil {
			return nil, err
		}
	}
	return d.run(c, args)
}

func truthy(v any) bool {
	switch t := v.(type) {
	case bool:
		return t
	case string:
		s := strings.ToLower(strings.TrimSpace(t))
		return s == "1" || s == "true" || s == "yes" || s == "y"
	case float64:
		return t != 0
	default:
		return false
	}
}

func strArg(args map[string]any, key string) string {
	v, ok := args[key]
	if !ok || v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case float64:
		// JSON numbers — 大雪花 ID 可能丢精度；尽量用字符串传
		return fmt.Sprintf("%.0f", t)
	case json.Number:
		return t.String()
	default:
		b, _ := json.Marshal(t)
		return strings.Trim(string(b), `"`)
	}
}

func requireStr(args map[string]any, keys ...string) (map[string]string, error) {
	out := map[string]string{}
	for _, k := range keys {
		s := strings.TrimSpace(strArg(args, k))
		if s == "" {
			return nil, fmt.Errorf("缺少参数 %s", k)
		}
		out[k] = s
	}
	return out, nil
}

func optStr(args map[string]any, key string) string {
	return strings.TrimSpace(strArg(args, key))
}

func bodyFromArgs(args map[string]any, skip ...string) map[string]any {
	skipSet := map[string]struct{}{"confirm": {}}
	for _, s := range skip {
		skipSet[s] = struct{}{}
	}
	out := map[string]any{}
	for k, v := range args {
		if _, ok := skipSet[k]; ok {
			continue
		}
		if v == nil {
			continue
		}
		if s, ok := v.(string); ok && s == "" {
			continue
		}
		out[k] = v
	}
	return out
}

func schemaObject(props map[string]any, required ...string) map[string]any {
	s := map[string]any{
		"type":       "object",
		"properties": props,
	}
	if len(required) > 0 {
		s["required"] = required
	}
	return s
}

func propString(desc string) map[string]any {
	return map[string]any{"type": "string", "description": desc}
}

func propBool(desc string) map[string]any {
	return map[string]any{"type": "boolean", "description": desc}
}

func propNumber(desc string) map[string]any {
	return map[string]any{"type": "number", "description": desc}
}

func propInteger(desc string) map[string]any {
	return map[string]any{"type": "integer", "description": desc}
}
