// Package mcp 实现 Model Context Protocol（stdio JSON-RPC 2.0），
// 暴露 tools / resources / prompts 给 AI 宿主。
package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/NewtSpeak/Newt-Agent/internal/auth"
	"github.com/NewtSpeak/Newt-Agent/internal/tools"
)

const protocolVersion = "2024-11-05"
const serverName = "newt-agent"
const serverVersion = "0.4.0"

type Server struct {
	Reg    *tools.Registry
	Reader io.Reader
	Writer io.Writer
	ErrLog io.Writer
}

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type rpcResponse struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      any       `json:"id,omitempty"`
	Result  any       `json:"result,omitempty"`
	Error   *rpcError `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Serve 阻塞读取 stdin，写 stdout（每行一条 JSON-RPC）。
func (s *Server) Serve() error {
	if s.Reader == nil {
		s.Reader = os.Stdin
	}
	if s.Writer == nil {
		s.Writer = os.Stdout
	}
	if s.ErrLog == nil {
		s.ErrLog = os.Stderr
	}

	var writeMu sync.Mutex
	write := func(v any) error {
		writeMu.Lock()
		defer writeMu.Unlock()
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		_, err = s.Writer.Write(append(b, '\n'))
		return err
	}

	sc := bufio.NewScanner(s.Reader)
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, 8*1024*1024)

	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var req rpcRequest
		if err := json.Unmarshal(line, &req); err != nil {
			_ = write(rpcResponse{
				JSONRPC: "2.0",
				Error:   &rpcError{Code: -32700, Message: "Parse error"},
			})
			continue
		}
		isNotify := req.ID == nil || string(req.ID) == "null"
		result, err := s.dispatch(req.Method, req.Params)
		if isNotify {
			continue
		}
		var id any
		_ = json.Unmarshal(req.ID, &id)
		if err != nil {
			_ = write(rpcResponse{
				JSONRPC: "2.0", ID: id,
				Error: &rpcError{Code: -32000, Message: err.Error()},
			})
			continue
		}
		if err := write(rpcResponse{JSONRPC: "2.0", ID: id, Result: result}); err != nil {
			return err
		}
	}
	return sc.Err()
}

func (s *Server) dispatch(method string, params json.RawMessage) (any, error) {
	switch method {
	case "initialize":
		return map[string]any{
			"protocolVersion": protocolVersion,
			"capabilities": map[string]any{
				"tools":     map[string]any{},
				"resources": map[string]any{},
				"prompts":   map[string]any{},
			},
			"serverInfo": map[string]any{
				"name":    serverName,
				"version": serverVersion,
			},
		}, nil
	case "notifications/initialized", "initialized":
		return nil, nil
	case "ping":
		return map[string]any{}, nil

	case "tools/list":
		defs := s.Reg.List()
		toolsOut := make([]map[string]any, 0, len(defs))
		for _, d := range defs {
			toolsOut = append(toolsOut, map[string]any{
				"name":        d.Name,
				"description": d.Description,
				"inputSchema": d.InputSchema,
			})
		}
		return map[string]any{"tools": toolsOut}, nil
	case "tools/call":
		var p struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("invalid tools/call params: %w", err)
		}
		result, err := s.Reg.Call(p.Name, p.Arguments, false)
		if err != nil {
			return map[string]any{
				"content": []map[string]any{{"type": "text", "text": "Error: " + err.Error()}},
				"isError": true,
			}, nil
		}
		text, ferr := tools.FormatResult(result)
		if ferr != nil {
			text = fmt.Sprint(result)
		}
		return map[string]any{
			"content": []map[string]any{{"type": "text", "text": text}},
			"isError": false,
		}, nil

	case "resources/list":
		return map[string]any{"resources": listResources()}, nil
	case "resources/read":
		var p struct {
			URI string `json:"uri"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("invalid resources/read params: %w", err)
		}
		return readResource(s, p.URI)
	case "resources/templates/list":
		return map[string]any{"resourceTemplates": []any{}}, nil

	case "prompts/list":
		return map[string]any{"prompts": listPrompts()}, nil
	case "prompts/get":
		var p struct {
			Name      string            `json:"name"`
			Arguments map[string]string `json:"arguments"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("invalid prompts/get params: %w", err)
		}
		return getPrompt(p.Name, p.Arguments)

	default:
		return nil, fmt.Errorf("method not found: %s", method)
	}
}

func listResources() []map[string]any {
	return []map[string]any{
		{
			"uri":         "newtspeak://status",
			"name":        "session-status",
			"description": "当前 newt profile 登录状态与 token 存储后端",
			"mimeType":    "application/json",
		},
		{
			"uri":         "newtspeak://tools",
			"name":        "tool-catalog",
			"description": "全部 tool 名称与说明（只读目录）",
			"mimeType":    "application/json",
		},
		{
			"uri":         "newtspeak://whoami",
			"name":        "current-user",
			"description": "OAuth 用户信息（需已登录）",
			"mimeType":    "application/json",
		},
		{
			"uri":         "newtspeak://guilds",
			"name":        "my-guilds",
			"description": "已加入的服务器列表（需已登录）",
			"mimeType":    "application/json",
		},
	}
}

func readResource(s *Server, uri string) (any, error) {
	uri = strings.TrimSpace(uri)
	var text string
	var err error
	switch uri {
	case "newtspeak://status":
		text, err = tools.FormatResult(auth.SessionMeta())
	case "newtspeak://tools":
		text, err = tools.FormatResult(s.Reg.List())
	case "newtspeak://whoami":
		res, e := s.Reg.Call("whoami", nil, false)
		if e != nil {
			return nil, e
		}
		text, err = tools.FormatResult(res)
	case "newtspeak://guilds":
		res, e := s.Reg.Call("guilds.list", nil, false)
		if e != nil {
			return nil, e
		}
		text, err = tools.FormatResult(res)
	default:
		return nil, fmt.Errorf("unknown resource uri: %s", uri)
	}
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"contents": []map[string]any{
			{"uri": uri, "mimeType": "application/json", "text": text},
		},
	}, nil
}

func listPrompts() []map[string]any {
	return []map[string]any{
		{
			"name":        "moderate-guild",
			"description": "协助管理某个服务器：先列频道/成员，再按用户意图执行治理",
			"arguments": []map[string]any{
				{"name": "guild_id", "description": "服务器 ID", "required": true},
				{"name": "goal", "description": "管理目标（如：整理频道、处理违规）", "required": false},
			},
		},
		{
			"name":        "audit-review",
			"description": "审查服务器审计日志并总结近期操作",
			"arguments": []map[string]any{
				{"name": "guild_id", "description": "服务器 ID", "required": true},
				{"name": "action_prefix", "description": "action 前缀过滤", "required": false},
			},
		},
		{
			"name":        "safe-ops",
			"description": "提醒 AI 使用 owl 工具时的安全约束（confirm、禁止索要密码）",
			"arguments":   []map[string]any{},
		},
	}
}

func getPrompt(name string, args map[string]string) (any, error) {
	if args == nil {
		args = map[string]string{}
	}
	var text string
	switch name {
	case "moderate-guild":
		gid := args["guild_id"]
		goal := args["goal"]
		if goal == "" {
			goal = "日常治理"
		}
		text = fmt.Sprintf(`你正在通过 NewtSpeak MCP/CLI 管理服务器。

服务器 ID: %s
目标: %s

建议步骤:
1. 调用 guilds.get / channels.list / members.list 了解现状
2. 向用户确认拟执行的写操作
3. 危险工具必须传 confirm=true
4. 不要索要用户密码或 refresh token
5. 完成后用 audit.list 核对（若有权限）`, gid, goal)
	case "audit-review":
		gid := args["guild_id"]
		prefix := args["action_prefix"]
		text = fmt.Sprintf(`审查 NewtSpeak 服务器审计日志。

guild_id: %s
action 前缀过滤: %s（可空）

请调用 audit.list，总结：谁做了什么、是否异常、是否需要后续限制/踢封。`, gid, prefix)
	case "safe-ops":
		text = `NewtSpeak Agent 安全规则:
- 身份是用户委托 OAuth（aud=agent），不是 Bot
- 禁止要求用户提供密码、refresh token、access token
- 踢人/删频道/封禁/平台写操作等 destructive 工具必须 confirm=true，并先征得用户明确同意
- 权限受用户 RBAC 限制；平台 API 还需 system_admin + platform scope
- 优先只读摸底，再提案，最后执行`
	default:
		return nil, fmt.Errorf("unknown prompt: %s", name)
	}
	return map[string]any{
		"description": name,
		"messages": []map[string]any{
			{
				"role": "user",
				"content": map[string]any{
					"type": "text",
					"text": text,
				},
			},
		},
	}, nil
}
