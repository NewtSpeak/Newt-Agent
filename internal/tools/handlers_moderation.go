package tools

import (
	"fmt"
	"net/http"

	"github.com/OwlSpeak/Owl-Agent/internal/api"
)

func (r *Registry) registerModerationExtras() {
	// ---- Restriction ----
	r.add(&Def{
		Name: "restrictions.list", Description: "列出服务器限制（需 MODERATE_MEMBERS）",
		InputSchema: schemaObject(map[string]any{
			"guild_id":   propString("服务器 ID"),
			"user_id":    propString("按目标用户过滤"),
			"channel_id": propString("按频道过滤"),
			"scope":      propString("范围过滤"),
			"active":     propString("true|false 仅活跃/已解除"),
		}, "guild_id"),
		CLIHint: "owl restrictions list --guild ...",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id")
			if err != nil {
				return nil, err
			}
			q := map[string]string{}
			for _, k := range []string{"user_id", "channel_id", "scope", "active"} {
				if v := optStr(args, k); v != "" {
					q[k] = v
				}
			}
			return c.Gapi(http.MethodGet, "/guilds/"+req["guild_id"]+"/restrictions", nil, q)
		},
	})
	r.add(&Def{
		Name: "restrictions.create", Description: "创建限制（禁言/禁视等）。scope: TEXT_CHANNEL|VOICE_CHANNEL|GUILD_ALL_TEXT|GUILD_ALL_VOICE；kind: SANCTION|CHANNEL_BAN",
		Destructive: true,
		InputSchema: schemaObject(map[string]any{
			"guild_id":       propString("服务器 ID"),
			"target_user_id": propString("目标用户 ID"),
			"scope":          propString("限制范围"),
			"kind":           propString("SANCTION 或 CHANNEL_BAN"),
			"channel_id":     propString("频道级限制时必填"),
			"reason":         propString("原因"),
			"expires_at":     propString("RFC3339 过期时间；省略=永久"),
			"deny": map[string]any{
				"type": "object",
				"description": "拒绝标志 view_text/send_text/listen_voice/speak_voice",
				"properties": map[string]any{
					"view_text":    propBool("禁止查看文字"),
					"send_text":    propBool("禁止发送文字"),
					"listen_voice": propBool("禁止听语音"),
					"speak_voice":  propBool("禁止发言"),
				},
			},
			"confirm": propBool("必须为 true"),
		}, "guild_id", "target_user_id", "scope", "kind", "confirm"),
		CLIHint: "owl restrictions create --guild ... --user ... --scope GUILD_ALL_TEXT --kind SANCTION --yes",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id", "target_user_id", "scope", "kind")
			if err != nil {
				return nil, err
			}
			body := map[string]any{
				"target_user_id": req["target_user_id"],
				"scope":          req["scope"],
				"kind":           req["kind"],
			}
			if v := optStr(args, "channel_id"); v != "" {
				body["channel_id"] = v
			}
			if v := optStr(args, "reason"); v != "" {
				body["reason"] = v
			}
			if v := optStr(args, "expires_at"); v != "" {
				body["expires_at"] = v
			}
			if d, ok := args["deny"]; ok {
				body["deny"] = d
			} else {
				// 默认禁发文字
				body["deny"] = map[string]any{"send_text": true}
			}
			return c.Gapi(http.MethodPost, "/guilds/"+req["guild_id"]+"/restrictions", body, nil)
		},
	})
	r.add(&Def{
		Name: "restrictions.get", Description: "限制详情",
		InputSchema: schemaObject(map[string]any{
			"guild_id":       propString("服务器 ID"),
			"restriction_id": propString("限制 ID"),
		}, "guild_id", "restriction_id"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id", "restriction_id")
			if err != nil {
				return nil, err
			}
			path := fmt.Sprintf("/guilds/%s/restrictions/%s", req["guild_id"], req["restriction_id"])
			return c.Gapi(http.MethodGet, path, nil, nil)
		},
	})
	r.add(&Def{
		Name: "restrictions.patch", Description: "更新限制 reason/expires_at",
		InputSchema: schemaObject(map[string]any{
			"guild_id":       propString("服务器 ID"),
			"restriction_id": propString("限制 ID"),
			"reason":         propString("原因"),
			"expires_at":     propString("RFC3339；可空"),
		}, "guild_id", "restriction_id"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id", "restriction_id")
			if err != nil {
				return nil, err
			}
			body := bodyFromArgs(args, "guild_id", "restriction_id")
			path := fmt.Sprintf("/guilds/%s/restrictions/%s", req["guild_id"], req["restriction_id"])
			return c.Gapi(http.MethodPatch, path, body, nil)
		},
	})
	r.add(&Def{
		Name: "restrictions.lift", Description: "提前解除限制",
		Destructive: true,
		InputSchema: schemaObject(map[string]any{
			"guild_id":       propString("服务器 ID"),
			"restriction_id": propString("限制 ID"),
			"confirm":        propBool("必须为 true"),
		}, "guild_id", "restriction_id", "confirm"),
		CLIHint: "owl restrictions lift --guild ... --id ... --yes",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id", "restriction_id")
			if err != nil {
				return nil, err
			}
			path := fmt.Sprintf("/guilds/%s/restrictions/%s", req["guild_id"], req["restriction_id"])
			return c.Gapi(http.MethodDelete, path, nil, nil)
		},
	})

	// ---- 审计 ----
	r.add(&Def{
		Name: "audit.list", Description: "服务器审计日志（需 VIEW_AUDIT_LOG）",
		InputSchema: schemaObject(map[string]any{
			"guild_id":    propString("服务器 ID"),
			"actor_id":    propString("操作者"),
			"action":      propString("action 前缀，如 restriction. / moderation."),
			"target_type": propString("目标类型"),
			"since":       propString("起始时间"),
			"until":       propString("截止时间"),
			"limit":       propInteger("条数"),
			"before":      propString("游标"),
		}, "guild_id"),
		CLIHint: "owl audit list --guild ...",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id")
			if err != nil {
				return nil, err
			}
			q := map[string]string{}
			for _, k := range []string{"actor_id", "action", "target_type", "since", "until", "before"} {
				if v := optStr(args, k); v != "" {
					q[k] = v
				}
			}
			if v, ok := args["limit"]; ok {
				q["limit"] = fmt.Sprint(v)
			}
			return c.Gapi(http.MethodGet, "/guilds/"+req["guild_id"]+"/audit-logs", nil, q)
		},
	})

	// ---- 语音治理 ----
	r.add(&Def{
		Name: "voice.states", Description: "频道语音成员状态列表",
		InputSchema: schemaObject(map[string]any{
			"guild_id":   propString("服务器 ID"),
			"channel_id": propString("语音频道 ID"),
		}, "guild_id", "channel_id"),
		CLIHint: "owl voice states --guild ... --channel ...",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id", "channel_id")
			if err != nil {
				return nil, err
			}
			path := fmt.Sprintf("/guilds/%s/channels/%s/voice-states", req["guild_id"], req["channel_id"])
			return c.Gapi(http.MethodGet, path, nil, nil)
		},
	})
	r.add(&Def{
		Name: "voice.disconnect", Description: "将成员踢出语音（需 MOVE/MUTE 权限）",
		Destructive: true,
		InputSchema: schemaObject(map[string]any{
			"guild_id": propString("服务器 ID"),
			"user_id":  propString("目标用户 ID"),
			"confirm":  propBool("必须为 true"),
		}, "guild_id", "user_id", "confirm"),
		CLIHint: "owl voice disconnect --guild ... --user ... --yes",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id", "user_id")
			if err != nil {
				return nil, err
			}
			body := map[string]any{"user_id": req["user_id"]}
			return c.Gapi(http.MethodPost, "/guilds/"+req["guild_id"]+"/voice/disconnect", body, nil)
		},
	})
	r.add(&Def{
		Name: "voice.move", Description: "移动成员到另一语音频道",
		Destructive: true,
		InputSchema: schemaObject(map[string]any{
			"guild_id":   propString("服务器 ID"),
			"user_id":    propString("用户 ID"),
			"channel_id": propString("目标语音频道 ID"),
			"confirm":    propBool("必须为 true"),
		}, "guild_id", "user_id", "channel_id", "confirm"),
		CLIHint: "owl voice move --guild ... --user ... --channel ... --yes",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id", "user_id", "channel_id")
			if err != nil {
				return nil, err
			}
			body := map[string]any{
				"user_id":    req["user_id"],
				"channel_id": req["channel_id"],
			}
			return c.Gapi(http.MethodPost, "/guilds/"+req["guild_id"]+"/voice/move", body, nil)
		},
	})
	r.add(&Def{
		Name: "voice.server_mute", Description: "服务器静音/耳聋成员（server_mute / server_deaf）",
		Destructive: true,
		InputSchema: schemaObject(map[string]any{
			"guild_id":     propString("服务器 ID"),
			"user_id":      propString("用户 ID"),
			"server_mute":  propBool("服务器静音"),
			"server_deaf":  propBool("服务器耳聋"),
			"confirm":      propBool("必须为 true"),
		}, "guild_id", "user_id", "confirm"),
		CLIHint: "owl voice mute --guild ... --user ... --mute --yes",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id", "user_id")
			if err != nil {
				return nil, err
			}
			body := map[string]any{}
			if _, ok := args["server_mute"]; ok {
				body["server_mute"] = truthy(args["server_mute"])
			}
			if _, ok := args["server_deaf"]; ok {
				body["server_deaf"] = truthy(args["server_deaf"])
			}
			// CLI 简写 mute/deafen
			if _, ok := args["mute"]; ok {
				body["server_mute"] = truthy(args["mute"])
			}
			if _, ok := args["deafen"]; ok {
				body["server_deaf"] = truthy(args["deafen"])
			}
			if len(body) == 0 {
				return nil, fmt.Errorf("请提供 server_mute 和/或 server_deaf")
			}
			path := fmt.Sprintf("/guilds/%s/voice/states/%s", req["guild_id"], req["user_id"])
			return c.Gapi(http.MethodPatch, path, body, nil)
		},
	})
	r.add(&Def{
		Name: "voice.nodes", Description: "服务器语音候选节点（成员可读）",
		InputSchema: schemaObject(map[string]any{"guild_id": propString("服务器 ID")}, "guild_id"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id")
			if err != nil {
				return nil, err
			}
			return c.Gapi(http.MethodGet, "/guilds/"+req["guild_id"]+"/voice/nodes", nil, nil)
		},
	})
	r.add(&Def{
		Name: "voice.node_pool.get", Description: "获取服务器节点池配置（需 MANAGE_GUILD）",
		InputSchema: schemaObject(map[string]any{"guild_id": propString("服务器 ID")}, "guild_id"),
		CLIHint:     "owl voice node-pool get --guild ...",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id")
			if err != nil {
				return nil, err
			}
			return c.Gapi(http.MethodGet, "/guilds/"+req["guild_id"]+"/node-pool", nil, nil)
		},
	})
	r.add(&Def{
		Name: "voice.node_pool.set", Description: "设置服务器生效节点（从候选勾选）",
		InputSchema: schemaObject(map[string]any{
			"guild_id":            propString("服务器 ID"),
			"node_ids":            map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "节点 ID 列表"},
			"fallback_to_default": propBool("候选空时回退默认池"),
		}, "guild_id", "node_ids"),
		CLIHint: "owl voice node-pool set --guild ... --nodes id1,id2",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id")
			if err != nil {
				return nil, err
			}
			body := bodyFromArgs(args, "guild_id")
			return c.Gapi(http.MethodPut, "/guilds/"+req["guild_id"]+"/node-pool", body, nil)
		},
	})

	// ---- 平台管理（/api/v1，需 platform scope）----
	r.add(&Def{
		Name: "platform.users.list", Description: "平台用户目录（需 system_admin + platform scope）",
		InputSchema: schemaObject(map[string]any{
			"q":     propString("搜索"),
			"limit": propInteger("条数"),
			"offset": propInteger("偏移"),
		}),
		CLIHint: "owl platform users list",
		run: func(c *api.Client, args map[string]any) (any, error) {
			q := map[string]string{}
			if v := optStr(args, "q"); v != "" {
				q["q"] = v
			}
			if v, ok := args["limit"]; ok {
				q["limit"] = fmt.Sprint(v)
			}
			if v, ok := args["offset"]; ok {
				q["offset"] = fmt.Sprint(v)
			}
			return c.Api(http.MethodGet, "/admin/users", nil, q)
		},
	})
	r.add(&Def{
		Name: "platform.users.disable", Description: "禁用平台用户",
		Destructive: true,
		InputSchema: schemaObject(map[string]any{
			"user_id": propString("用户 ID"),
			"confirm": propBool("必须为 true"),
		}, "user_id", "confirm"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "user_id")
			if err != nil {
				return nil, err
			}
			return c.Api(http.MethodPost, "/admin/users/"+req["user_id"]+"/disable", map[string]any{}, nil)
		},
	})
	r.add(&Def{
		Name: "platform.users.enable", Description: "启用平台用户",
		InputSchema: schemaObject(map[string]any{"user_id": propString("用户 ID")}, "user_id"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "user_id")
			if err != nil {
				return nil, err
			}
			return c.Api(http.MethodPost, "/admin/users/"+req["user_id"]+"/enable", map[string]any{}, nil)
		},
	})
	r.add(&Def{
		Name: "platform.users.reset_password", Description: "管理员重置用户密码",
		Destructive: true,
		InputSchema: schemaObject(map[string]any{
			"user_id":      propString("用户 ID"),
			"new_password": propString("新密码 8-128"),
			"confirm":      propBool("必须为 true"),
		}, "user_id", "new_password", "confirm"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "user_id", "new_password")
			if err != nil {
				return nil, err
			}
			return c.Api(http.MethodPost, "/admin/users/"+req["user_id"]+"/reset-password", map[string]any{
				"new_password": req["new_password"],
			}, nil)
		},
	})
	r.add(&Def{
		Name: "platform.users.set_admin", Description: "授予/回收 system_admin",
		Destructive: true,
		InputSchema: schemaObject(map[string]any{
			"user_id":      propString("用户 ID"),
			"system_admin": propBool("是否系统管理员"),
			"confirm":      propBool("必须为 true"),
		}, "user_id", "confirm"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "user_id")
			if err != nil {
				return nil, err
			}
			body := map[string]any{"system_admin": truthy(args["system_admin"])}
			return c.Api(http.MethodPatch, "/admin/users/"+req["user_id"]+"/system-admin", body, nil)
		},
	})
	r.add(&Def{
		Name: "platform.registration.get", Description: "读取用户端注册开关",
		InputSchema: schemaObject(map[string]any{}),
		run: func(c *api.Client, _ map[string]any) (any, error) {
			return c.Api(http.MethodGet, "/admin/registration", nil, nil)
		},
	})
	r.add(&Def{
		Name: "platform.registration.set", Description: "设置用户端注册开关",
		Destructive: true,
		InputSchema: schemaObject(map[string]any{
			"signup_enabled": propBool("是否开放注册"),
			"confirm":        propBool("必须为 true"),
		}, "confirm"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			enabled := truthy(args["signup_enabled"])
			if v, ok := args["enabled"]; ok {
				enabled = truthy(v)
			}
			return c.Api(http.MethodPut, "/admin/registration", map[string]any{
				"signup_enabled": enabled,
			}, nil)
		},
	})
	r.add(&Def{
		Name: "platform.sfu.nodes", Description: "列出 SFU 节点（平台）",
		InputSchema: schemaObject(map[string]any{}),
		CLIHint:     "owl platform sfu nodes",
		run: func(c *api.Client, _ map[string]any) (any, error) {
			return c.Api(http.MethodGet, "/admin/sfu/nodes", nil, nil)
		},
	})
	r.add(&Def{
		Name: "platform.sfu.topology", Description: "SFU 拓扑",
		InputSchema: schemaObject(map[string]any{}),
		run: func(c *api.Client, _ map[string]any) (any, error) {
			return c.Api(http.MethodGet, "/admin/sfu/topology", nil, nil)
		},
	})
	r.add(&Def{
		Name: "platform.audit.list", Description: "全站审计日志（平台）",
		InputSchema: schemaObject(map[string]any{
			"actor_id": propString("操作者"),
			"action":   propString("action 前缀"),
			"limit":    propInteger("条数"),
			"before":   propString("游标"),
		}),
		run: func(c *api.Client, args map[string]any) (any, error) {
			q := map[string]string{}
			for _, k := range []string{"actor_id", "action", "before"} {
				if v := optStr(args, k); v != "" {
					q[k] = v
				}
			}
			if v, ok := args["limit"]; ok {
				q["limit"] = fmt.Sprint(v)
			}
			return c.Api(http.MethodGet, "/admin/audit-logs", nil, q)
		},
	})
}
