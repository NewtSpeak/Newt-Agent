package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/NewtSpeak/Newt-Agent/internal/api"
	"github.com/NewtSpeak/Newt-Agent/internal/auth"
	"github.com/NewtSpeak/Newt-Agent/internal/config"
)

func (r *Registry) registerAll() {
	// ---- 会话 / 账号 ----
	r.add(&Def{
		Name: "whoami", Description: "当前 OAuth 用户信息",
		InputSchema: schemaObject(map[string]any{}),
		CLIHint:     "newt whoami",
		run: func(c *api.Client, _ map[string]any) (any, error) {
			return c.UserInfo()
		},
	})
	r.add(&Def{
		Name: "status", Description: "CLI 登录状态与配置（含 token 存储后端）",
		InputSchema: schemaObject(map[string]any{}),
		CLIHint:     "newt status",
		run: func(_ *api.Client, _ map[string]any) (any, error) {
			return auth.SessionMeta(), nil
		},
	})

	// ---- 服务器 ----
	r.add(&Def{
		Name: "guilds.list", Description: "列出当前用户加入的服务器",
		InputSchema: schemaObject(map[string]any{}),
		CLIHint:     "newt guilds list",
		run: func(c *api.Client, _ map[string]any) (any, error) {
			return c.Gapi(http.MethodGet, "/users/@me/guilds", nil, nil)
		},
	})
	r.add(&Def{
		Name: "guilds.get", Description: "获取服务器详情",
		InputSchema: schemaObject(map[string]any{"guild_id": propString("服务器 ID")}, "guild_id"),
		CLIHint:     "newt guilds get <guild_id>",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id")
			if err != nil {
				return nil, err
			}
			return c.Gapi(http.MethodGet, "/guilds/"+req["guild_id"], nil, nil)
		},
	})
	r.add(&Def{
		Name: "guilds.create", Description: "创建服务器",
		InputSchema: schemaObject(map[string]any{"name": propString("服务器名称 2-100 字符")}, "name"),
		CLIHint:     "newt guilds create --name ...",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "name")
			if err != nil {
				return nil, err
			}
			return c.Gapi(http.MethodPost, "/guilds", map[string]any{"name": req["name"]}, nil)
		},
	})
	r.add(&Def{
		Name: "guilds.update", Description: "更新服务器（名称/简介/默认频道等，需 MANAGE_GUILD）",
		InputSchema: schemaObject(map[string]any{
			"guild_id":                   propString("服务器 ID"),
			"name":                       propString("新名称"),
			"description":                propString("简介"),
			"default_channel_id":         propString("默认着陆文字频道 ID；空串清空"),
			"restriction_badge_visible":  propBool("是否显示受限徽章"),
			"restriction_reason_required": propBool("Restriction 是否强制 reason"),
		}, "guild_id"),
		CLIHint: "newt guilds update --guild ... --name ...",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id")
			if err != nil {
				return nil, err
			}
			body := bodyFromArgs(args, "guild_id")
			return c.Gapi(http.MethodPatch, "/guilds/"+req["guild_id"], body, nil)
		},
	})
	r.add(&Def{
		Name: "guilds.delete", Description: "删除服务器（仅所有者；需 confirm_name 与名称一致）",
		Destructive: true,
		InputSchema: schemaObject(map[string]any{
			"guild_id":     propString("服务器 ID"),
			"confirm_name": propString("必须与当前服务器名称完全一致"),
			"confirm":      propBool("必须为 true"),
		}, "guild_id", "confirm_name", "confirm"),
		CLIHint: "newt guilds delete --guild ... --confirm-name ... --yes",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id", "confirm_name")
			if err != nil {
				return nil, err
			}
			return c.Gapi(http.MethodDelete, "/guilds/"+req["guild_id"], map[string]any{
				"confirm_name": req["confirm_name"],
			}, nil)
		},
	})
	r.add(&Def{
		Name: "guilds.permissions", Description: "查询本人在服务器的最终权限位",
		InputSchema: schemaObject(map[string]any{"guild_id": propString("服务器 ID")}, "guild_id"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id")
			if err != nil {
				return nil, err
			}
			return c.Gapi(http.MethodGet, "/guilds/"+req["guild_id"]+"/permissions/@me", nil, nil)
		},
	})
	r.add(&Def{
		Name: "guilds.transfer", Description: "转让服务器所有权（仅所有者）",
		Destructive: true,
		InputSchema: schemaObject(map[string]any{
			"guild_id":          propString("服务器 ID"),
			"new_owner_user_id": propString("新所有者用户 ID（须为本服成员）"),
			"confirm":           propBool("必须为 true"),
		}, "guild_id", "new_owner_user_id", "confirm"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id", "new_owner_user_id")
			if err != nil {
				return nil, err
			}
			return c.Gapi(http.MethodPost, "/guilds/"+req["guild_id"]+"/transfer-ownership", map[string]any{
				"new_owner_user_id": req["new_owner_user_id"],
			}, nil)
		},
	})

	// ---- 频道 ----
	r.add(&Def{
		Name: "channels.list", Description: "列出服务器可见频道",
		InputSchema: schemaObject(map[string]any{"guild_id": propString("服务器 ID")}, "guild_id"),
		CLIHint:     "newt channels list --guild ...",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id")
			if err != nil {
				return nil, err
			}
			return c.Gapi(http.MethodGet, "/guilds/"+req["guild_id"]+"/channels", nil, nil)
		},
	})
	r.add(&Def{
		Name: "channels.create", Description: "创建频道或分类（需 MANAGE_CHANNELS）。type: TEXT|VOICE|CATEGORY|STAGE",
		InputSchema: schemaObject(map[string]any{
			"guild_id":   propString("服务器 ID"),
			"name":       propString("频道名"),
			"type":       propString("TEXT | VOICE | CATEGORY | STAGE 等"),
			"topic":      propString("主题"),
			"parent_id":  propString("父分类 ID"),
			"position":   propInteger("排序位置"),
			"private":    propBool("私密频道"),
			"password":   propString("访问密码（上锁）"),
			"user_limit": propInteger("语音人数上限"),
		}, "guild_id", "name", "type"),
		CLIHint: "newt channels create --guild ... --name ... --type TEXT",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id", "name", "type")
			if err != nil {
				return nil, err
			}
			body := bodyFromArgs(args, "guild_id")
			body["name"] = req["name"]
			body["type"] = req["type"]
			return c.Gapi(http.MethodPost, "/guilds/"+req["guild_id"]+"/channels", body, nil)
		},
	})
	r.add(&Def{
		Name: "channels.update", Description: "更新频道（名称/主题/密码/父分类等，需 MANAGE_CHANNELS）",
		InputSchema: schemaObject(map[string]any{
			"channel_id": propString("频道 ID"),
			"name":       propString("新名称"),
			"topic":      propString("主题"),
			"parent_id":  propString("父分类；空/null 移出"),
			"password":   propString("设置密码"),
			"locked":     propBool("false 关锁"),
			"voice_note": propString("语音活动注释"),
			"user_limit": propInteger("语音上限"),
		}, "channel_id"),
		CLIHint: "newt channels update --channel ... --name ...",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "channel_id")
			if err != nil {
				return nil, err
			}
			body := bodyFromArgs(args, "channel_id")
			return c.Gapi(http.MethodPatch, "/channels/"+req["channel_id"], body, nil)
		},
	})
	r.add(&Def{
		Name: "channels.delete", Description: "删除频道或分类（需 MANAGE_CHANNELS）",
		Destructive: true,
		InputSchema: schemaObject(map[string]any{
			"channel_id": propString("频道 ID"),
			"confirm":    propBool("必须为 true"),
		}, "channel_id", "confirm"),
		CLIHint: "newt channels delete --channel ... --yes",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "channel_id")
			if err != nil {
				return nil, err
			}
			return c.Gapi(http.MethodDelete, "/channels/"+req["channel_id"], nil, nil)
		},
	})
	r.add(&Def{
		Name: "channels.reorder", Description: "批量排序/移动频道。items: [{id, position, parent_id?}]",
		InputSchema: schemaObject(map[string]any{
			"guild_id": propString("服务器 ID"),
			"items": map[string]any{
				"type": "array",
				"description": "排序项数组",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"id":        propString("频道 ID"),
						"position":  propInteger("位置"),
						"parent_id": propString("父分类，可 null"),
					},
				},
			},
		}, "guild_id", "items"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id")
			if err != nil {
				return nil, err
			}
			items, ok := args["items"]
			if !ok {
				return nil, fmt.Errorf("缺少 items")
			}
			return c.Gapi(http.MethodPatch, "/guilds/"+req["guild_id"]+"/channels", items, nil)
		},
	})
	r.add(&Def{
		Name: "channels.permissions", Description: "查询本人在频道的权限投影",
		InputSchema: schemaObject(map[string]any{"channel_id": propString("频道 ID")}, "channel_id"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "channel_id")
			if err != nil {
				return nil, err
			}
			return c.Gapi(http.MethodGet, "/channels/"+req["channel_id"]+"/permissions/@me", nil, nil)
		},
	})
	r.add(&Def{
		Name: "channels.overwrites.list", Description: "列出频道权限覆盖",
		InputSchema: schemaObject(map[string]any{
			"guild_id":   propString("服务器 ID"),
			"channel_id": propString("频道 ID"),
		}, "guild_id", "channel_id"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id", "channel_id")
			if err != nil {
				return nil, err
			}
			path := fmt.Sprintf("/guilds/%s/channels/%s/overwrites", req["guild_id"], req["channel_id"])
			return c.Gapi(http.MethodGet, path, nil, nil)
		},
	})
	r.add(&Def{
		Name: "channels.overwrites.upsert", Description: "创建/更新频道权限覆盖（需 MANAGE_ROLES）",
		InputSchema: schemaObject(map[string]any{
			"channel_id": propString("频道 ID"),
			"target_id":  propString("角色 ID 或成员记录 ID"),
			"type":       propString("ROLE 或 MEMBER"),
			"allow":      propNumber("允许位掩码"),
			"deny":       propNumber("拒绝位掩码"),
		}, "channel_id", "target_id", "type"),
		CLIHint: "newt channels overwrites set --channel ... --target ... --type ROLE --allow 0 --deny 0",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "channel_id", "target_id", "type")
			if err != nil {
				return nil, err
			}
			body := map[string]any{"type": req["type"]}
			if v, ok := args["allow"]; ok {
				body["allow"] = v
			} else {
				body["allow"] = 0
			}
			if v, ok := args["deny"]; ok {
				body["deny"] = v
			} else {
				body["deny"] = 0
			}
			path := fmt.Sprintf("/channels/%s/overwrites/%s", req["channel_id"], req["target_id"])
			return c.Gapi(http.MethodPut, path, body, nil)
		},
	})
	r.add(&Def{
		Name: "channels.overwrites.delete", Description: "删除频道权限覆盖",
		Destructive: true,
		InputSchema: schemaObject(map[string]any{
			"channel_id": propString("频道 ID"),
			"target_id":  propString("目标 ID"),
			"type":       propString("可选 ROLE|MEMBER 查询参数"),
			"confirm":    propBool("必须为 true"),
		}, "channel_id", "target_id", "confirm"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "channel_id", "target_id")
			if err != nil {
				return nil, err
			}
			q := map[string]string{}
			if t := optStr(args, "type"); t != "" {
				q["type"] = t
			}
			path := fmt.Sprintf("/channels/%s/overwrites/%s", req["channel_id"], req["target_id"])
			return c.Gapi(http.MethodDelete, path, nil, q)
		},
	})

	// ---- 角色 ----
	r.add(&Def{
		Name: "roles.list", Description: "列出服务器角色",
		InputSchema: schemaObject(map[string]any{"guild_id": propString("服务器 ID")}, "guild_id"),
		CLIHint:     "owl roles list --guild ...",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id")
			if err != nil {
				return nil, err
			}
			return c.Gapi(http.MethodGet, "/guilds/"+req["guild_id"]+"/roles", nil, nil)
		},
	})
	r.add(&Def{
		Name: "roles.create", Description: "创建角色（需 MANAGE_ROLES）",
		InputSchema: schemaObject(map[string]any{
			"guild_id":    propString("服务器 ID"),
			"name":        propString("角色名"),
			"permissions": propNumber("权限位掩码"),
			"position":    propInteger("位置 ≥1"),
			"color":       propString("颜色"),
			"hoist":       propBool("是否单独显示"),
			"mentionable": propBool("是否可提及"),
		}, "guild_id", "name"),
		CLIHint: "owl roles create --guild ... --name ... --permissions 0",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id", "name")
			if err != nil {
				return nil, err
			}
			body := bodyFromArgs(args, "guild_id")
			body["name"] = req["name"]
			if _, ok := body["permissions"]; !ok {
				body["permissions"] = 0
			}
			if _, ok := body["position"]; !ok {
				body["position"] = 1
			}
			return c.Gapi(http.MethodPost, "/guilds/"+req["guild_id"]+"/roles", body, nil)
		},
	})
	r.add(&Def{
		Name: "roles.update", Description: "更新角色",
		InputSchema: schemaObject(map[string]any{
			"guild_id":    propString("服务器 ID"),
			"role_id":     propString("角色 ID"),
			"name":        propString("角色名"),
			"permissions": propNumber("权限位"),
			"position":    propInteger("位置"),
			"color":       propString("颜色"),
			"hoist":       propBool("单独显示"),
			"mentionable": propBool("可提及"),
		}, "guild_id", "role_id", "name"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id", "role_id", "name")
			if err != nil {
				return nil, err
			}
			body := bodyFromArgs(args, "guild_id", "role_id")
			body["name"] = req["name"]
			if _, ok := body["permissions"]; !ok {
				body["permissions"] = 0
			}
			if _, ok := body["position"]; !ok {
				body["position"] = 1
			}
			path := fmt.Sprintf("/guilds/%s/roles/%s", req["guild_id"], req["role_id"])
			return c.Gapi(http.MethodPatch, path, body, nil)
		},
	})
	r.add(&Def{
		Name: "roles.delete", Description: "删除角色",
		Destructive: true,
		InputSchema: schemaObject(map[string]any{
			"guild_id": propString("服务器 ID"),
			"role_id":  propString("角色 ID"),
			"confirm":  propBool("必须为 true"),
		}, "guild_id", "role_id", "confirm"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id", "role_id")
			if err != nil {
				return nil, err
			}
			path := fmt.Sprintf("/guilds/%s/roles/%s", req["guild_id"], req["role_id"])
			return c.Gapi(http.MethodDelete, path, nil, nil)
		},
	})
	r.add(&Def{
		Name: "roles.assign", Description: "给成员赋予角色",
		InputSchema: schemaObject(map[string]any{
			"guild_id":  propString("服务器 ID"),
			"member_id": propString("成员 ID 或用户 ID"),
			"role_id":   propString("角色 ID"),
		}, "guild_id", "member_id", "role_id"),
		CLIHint: "owl roles assign --guild ... --member ... --role ...",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id", "member_id", "role_id")
			if err != nil {
				return nil, err
			}
			path := fmt.Sprintf("/guilds/%s/members/%s/roles/%s", req["guild_id"], req["member_id"], req["role_id"])
			return c.Gapi(http.MethodPut, path, nil, nil)
		},
	})
	r.add(&Def{
		Name: "roles.remove", Description: "移除成员角色",
		Destructive: true,
		InputSchema: schemaObject(map[string]any{
			"guild_id":  propString("服务器 ID"),
			"member_id": propString("成员 ID"),
			"role_id":   propString("角色 ID"),
			"confirm":   propBool("必须为 true"),
		}, "guild_id", "member_id", "role_id", "confirm"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id", "member_id", "role_id")
			if err != nil {
				return nil, err
			}
			path := fmt.Sprintf("/guilds/%s/members/%s/roles/%s", req["guild_id"], req["member_id"], req["role_id"])
			return c.Gapi(http.MethodDelete, path, nil, nil)
		},
	})

	// ---- 成员 / 治理 ----
	r.add(&Def{
		Name: "members.list", Description: "列出服务器成员",
		InputSchema: schemaObject(map[string]any{"guild_id": propString("服务器 ID")}, "guild_id"),
		CLIHint:     "owl members list --guild ...",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id")
			if err != nil {
				return nil, err
			}
			return c.Gapi(http.MethodGet, "/guilds/"+req["guild_id"]+"/members", nil, nil)
		},
	})
	r.add(&Def{
		Name: "members.nick", Description: "修改成员昵称",
		InputSchema: schemaObject(map[string]any{
			"guild_id":  propString("服务器 ID"),
			"member_id": propString("成员 ID 或 @me"),
			"nick":      propString("昵称；空串清除"),
		}, "guild_id", "member_id"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id", "member_id")
			if err != nil {
				return nil, err
			}
			nick := optStr(args, "nick")
			path := fmt.Sprintf("/guilds/%s/members/%s", req["guild_id"], req["member_id"])
			return c.Gapi(http.MethodPatch, path, map[string]any{"nick": nick}, nil)
		},
	})
	r.add(&Def{
		Name: "members.kick", Description: "踢出成员（@me 表示自己退服）",
		Destructive: true,
		InputSchema: schemaObject(map[string]any{
			"guild_id":  propString("服务器 ID"),
			"member_id": propString("成员 ID 或 @me"),
			"confirm":   propBool("必须为 true"),
		}, "guild_id", "member_id", "confirm"),
		CLIHint: "owl members kick --guild ... --member ... --yes",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id", "member_id")
			if err != nil {
				return nil, err
			}
			path := fmt.Sprintf("/guilds/%s/members/%s", req["guild_id"], req["member_id"])
			return c.Gapi(http.MethodDelete, path, nil, nil)
		},
	})
	r.add(&Def{
		Name: "members.ban", Description: "封禁用户",
		Destructive: true,
		InputSchema: schemaObject(map[string]any{
			"guild_id": propString("服务器 ID"),
			"user_id":  propString("用户 ID"),
			"reason":   propString("原因"),
			"confirm":  propBool("必须为 true"),
		}, "guild_id", "user_id", "confirm"),
		CLIHint: "owl members ban --guild ... --user ... --yes",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id", "user_id")
			if err != nil {
				return nil, err
			}
			body := map[string]any{}
			if r := optStr(args, "reason"); r != "" {
				body["reason"] = r
			}
			path := fmt.Sprintf("/guilds/%s/bans/%s", req["guild_id"], req["user_id"])
			return c.Gapi(http.MethodPut, path, body, nil)
		},
	})
	r.add(&Def{
		Name: "members.unban", Description: "解除封禁",
		InputSchema: schemaObject(map[string]any{
			"guild_id": propString("服务器 ID"),
			"user_id":  propString("用户 ID"),
		}, "guild_id", "user_id"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id", "user_id")
			if err != nil {
				return nil, err
			}
			path := fmt.Sprintf("/guilds/%s/bans/%s", req["guild_id"], req["user_id"])
			return c.Gapi(http.MethodDelete, path, nil, nil)
		},
	})
	r.add(&Def{
		Name: "members.bans", Description: "封禁列表",
		InputSchema: schemaObject(map[string]any{"guild_id": propString("服务器 ID")}, "guild_id"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id")
			if err != nil {
				return nil, err
			}
			return c.Gapi(http.MethodGet, "/guilds/"+req["guild_id"]+"/bans", nil, nil)
		},
	})

	// ---- 邀请 ----
	r.add(&Def{
		Name: "invites.create", Description: "创建服务器邀请",
		InputSchema: schemaObject(map[string]any{
			"guild_id":    propString("服务器 ID"),
			"ttl_seconds": propInteger("有效秒数，≥60；省略不限"),
			"max_uses":    propInteger("最大使用次数，≥1；省略不限"),
		}, "guild_id"),
		CLIHint: "owl invites create --guild ...",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id")
			if err != nil {
				return nil, err
			}
			body := map[string]any{}
			if v, ok := args["ttl_seconds"]; ok {
				body["ttl_seconds"] = v
			}
			if v, ok := args["max_uses"]; ok {
				body["max_uses"] = v
			}
			return c.Gapi(http.MethodPost, "/guilds/"+req["guild_id"]+"/invites", body, nil)
		},
	})
	r.add(&Def{
		Name: "invites.get", Description: "预览邀请码信息",
		InputSchema: schemaObject(map[string]any{"code": propString("邀请码")}, "code"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "code")
			if err != nil {
				return nil, err
			}
			return c.Gapi(http.MethodGet, "/invites/"+req["code"], nil, nil)
		},
	})
	r.add(&Def{
		Name: "invites.join", Description: "凭邀请码加入服务器",
		InputSchema: schemaObject(map[string]any{"code": propString("邀请码")}, "code"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "code")
			if err != nil {
				return nil, err
			}
			return c.Gapi(http.MethodPost, "/invites/"+req["code"]+"/join", nil, nil)
		},
	})

	// ---- 消息 ----
	r.add(&Def{
		Name: "messages.list", Description: "拉取频道消息历史",
		InputSchema: schemaObject(map[string]any{
			"channel_id": propString("频道 ID"),
			"before":     propString("向前翻页消息 ID"),
			"after":      propString("向后翻页消息 ID"),
			"limit":      propInteger("条数，默认服务端默认"),
		}, "channel_id"),
		CLIHint: "owl messages list --channel ...",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "channel_id")
			if err != nil {
				return nil, err
			}
			q := map[string]string{}
			if v := optStr(args, "before"); v != "" {
				q["before"] = v
			}
			if v := optStr(args, "after"); v != "" {
				q["after"] = v
			}
			if v, ok := args["limit"]; ok {
				q["limit"] = fmt.Sprint(v)
			}
			return c.Gapi(http.MethodGet, "/channels/"+req["channel_id"]+"/messages", nil, q)
		},
	})
	r.add(&Def{
		Name: "messages.get", Description: "获取单条消息",
		InputSchema: schemaObject(map[string]any{
			"channel_id": propString("频道 ID"),
			"message_id": propString("消息 ID"),
		}, "channel_id", "message_id"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "channel_id", "message_id")
			if err != nil {
				return nil, err
			}
			path := fmt.Sprintf("/channels/%s/messages/%s", req["channel_id"], req["message_id"])
			return c.Gapi(http.MethodGet, path, nil, nil)
		},
	})
	r.add(&Def{
		Name: "messages.send", Description: "发送文本消息",
		InputSchema: schemaObject(map[string]any{
			"channel_id": propString("频道 ID"),
			"content":    propString("消息正文"),
		}, "channel_id", "content"),
		CLIHint: "owl messages send --channel ... --content ...",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "channel_id", "content")
			if err != nil {
				return nil, err
			}
			return c.Gapi(http.MethodPost, "/channels/"+req["channel_id"]+"/messages", map[string]any{
				"content": req["content"],
			}, nil)
		},
	})
	r.add(&Def{
		Name: "messages.edit", Description: "编辑消息（通常仅作者）",
		InputSchema: schemaObject(map[string]any{
			"channel_id": propString("频道 ID"),
			"message_id": propString("消息 ID"),
			"content":    propString("新正文"),
		}, "channel_id", "message_id", "content"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "channel_id", "message_id", "content")
			if err != nil {
				return nil, err
			}
			path := fmt.Sprintf("/channels/%s/messages/%s", req["channel_id"], req["message_id"])
			return c.Gapi(http.MethodPatch, path, map[string]any{"content": req["content"]}, nil)
		},
	})
	r.add(&Def{
		Name: "messages.delete", Description: "删除消息",
		Destructive: true,
		InputSchema: schemaObject(map[string]any{
			"channel_id": propString("频道 ID"),
			"message_id": propString("消息 ID"),
			"confirm":    propBool("必须为 true"),
		}, "channel_id", "message_id", "confirm"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "channel_id", "message_id")
			if err != nil {
				return nil, err
			}
			path := fmt.Sprintf("/channels/%s/messages/%s", req["channel_id"], req["message_id"])
			return c.Gapi(http.MethodDelete, path, nil, nil)
		},
	})
	r.add(&Def{
		Name: "messages.search", Description: "搜索消息（支持 guild/channel/author/before/after/limit）",
		InputSchema: schemaObject(map[string]any{
			"q":          propString("搜索关键词"),
			"guild_id":   propString("限定服务器"),
			"channel_id": propString("限定频道"),
			"author_id":  propString("限定作者用户 ID"),
			"before":     propString("消息 ID 游标（更旧）"),
			"after":      propString("消息 ID 游标（更新）"),
			"limit":      propInteger("1-100"),
		}, "q"),
		CLIHint: "owl messages search --q ... [--guild] [--channel] [--author] [--limit]",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "q")
			if err != nil {
				return nil, err
			}
			q := map[string]string{"q": req["q"]}
			for _, k := range []string{"guild_id", "channel_id", "author_id", "before", "after"} {
				if v := optStr(args, k); v != "" {
					q[k] = v
				}
			}
			if v, ok := args["limit"]; ok {
				q["limit"] = fmt.Sprint(v)
			}
			return c.Gapi(http.MethodGet, "/search/messages", nil, q)
		},
	})

	// 确保 strconv 可用（部分路径可能用）
	_ = strconv.Itoa
	_ = config.Load

	// Restriction / 审计 / 语音 / 平台
	r.registerModerationExtras()
	// 贴图 / 社交
	r.registerSocialAndStickers()
}

// FormatResult 将结果格式化为可读 JSON 文本（MCP content）。
func FormatResult(v any) (string, error) {
	switch t := v.(type) {
	case json.RawMessage:
		var pretty any
		if err := json.Unmarshal(t, &pretty); err != nil {
			return string(t), nil
		}
		b, err := json.MarshalIndent(pretty, "", "  ")
		return string(b), err
	case []byte:
		return string(t), nil
	case string:
		return t, nil
	default:
		b, err := json.MarshalIndent(t, "", "  ")
		return string(b), err
	}
}
