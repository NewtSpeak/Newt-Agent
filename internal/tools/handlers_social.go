package tools

import (
	"fmt"
	"net/http"

	"github.com/NewtSpeak/Newt-Agent/internal/api"
)

func (r *Registry) registerSocialAndStickers() {
	// ---- 贴图 ----
	r.add(&Def{
		Name: "stickers.packs.list", Description: "列出我的贴图包",
		InputSchema: schemaObject(map[string]any{}),
		CLIHint:     "owl stickers packs list",
		run: func(c *api.Client, _ map[string]any) (any, error) {
			return c.Gapi(http.MethodGet, "/users/@me/sticker-packs", nil, nil)
		},
	})
	r.add(&Def{
		Name: "stickers.packs.get", Description: "贴图包详情",
		InputSchema: schemaObject(map[string]any{"pack_id": propString("贴图包 ID")}, "pack_id"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "pack_id")
			if err != nil {
				return nil, err
			}
			return c.Gapi(http.MethodGet, "/sticker-packs/"+req["pack_id"], nil, nil)
		},
	})
	r.add(&Def{
		Name: "stickers.packs.create", Description: "创建贴图包",
		InputSchema: schemaObject(map[string]any{
			"name":        propString("名称"),
			"description": propString("描述"),
			"kind":        propString("类型，如 static/animated"),
		}, "name"),
		CLIHint: "owl stickers packs create --name ...",
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "name")
			if err != nil {
				return nil, err
			}
			body := bodyFromArgs(args)
			body["name"] = req["name"]
			return c.Gapi(http.MethodPost, "/users/@me/sticker-packs", body, nil)
		},
	})
	r.add(&Def{
		Name: "stickers.packs.delete", Description: "软删除贴图包",
		Destructive: true,
		InputSchema: schemaObject(map[string]any{
			"pack_id": propString("贴图包 ID"),
			"confirm": propBool("必须为 true"),
		}, "pack_id", "confirm"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "pack_id")
			if err != nil {
				return nil, err
			}
			return c.Gapi(http.MethodDelete, "/users/@me/sticker-packs/"+req["pack_id"], nil, nil)
		},
	})
	r.add(&Def{
		Name: "stickers.library.list", Description: "我的贴图库（已安装包）",
		InputSchema: schemaObject(map[string]any{}),
		CLIHint:     "owl stickers library list",
		run: func(c *api.Client, _ map[string]any) (any, error) {
			return c.Gapi(http.MethodGet, "/users/@me/sticker-library", nil, nil)
		},
	})
	r.add(&Def{
		Name: "stickers.library.install", Description: "安装贴图包到库",
		InputSchema: schemaObject(map[string]any{"pack_id": propString("贴图包 ID")}, "pack_id"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "pack_id")
			if err != nil {
				return nil, err
			}
			return c.Gapi(http.MethodPut, "/users/@me/sticker-library/"+req["pack_id"], nil, nil)
		},
	})
	r.add(&Def{
		Name: "stickers.library.uninstall", Description: "从库卸载贴图包",
		Destructive: true,
		InputSchema: schemaObject(map[string]any{
			"pack_id": propString("贴图包 ID"),
			"confirm": propBool("必须为 true"),
		}, "pack_id", "confirm"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "pack_id")
			if err != nil {
				return nil, err
			}
			return c.Gapi(http.MethodDelete, "/users/@me/sticker-library/"+req["pack_id"], nil, nil)
		},
	})
	r.add(&Def{
		Name: "stickers.available", Description: "可用贴图集合（可选 guild_id 过滤 ban）",
		InputSchema: schemaObject(map[string]any{"guild_id": propString("服务器 ID 过滤")}),
		CLIHint:     "owl stickers available [--guild]",
		run: func(c *api.Client, args map[string]any) (any, error) {
			q := map[string]string{}
			if v := optStr(args, "guild_id"); v != "" {
				q["guild_id"] = v
			}
			return c.Gapi(http.MethodGet, "/users/@me/sticker-available", nil, q)
		},
	})
	r.add(&Def{
		Name: "stickers.guild_bans.list", Description: "服务器贴图包 ban 列表",
		InputSchema: schemaObject(map[string]any{"guild_id": propString("服务器 ID")}, "guild_id"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id")
			if err != nil {
				return nil, err
			}
			return c.Gapi(http.MethodGet, "/guilds/"+req["guild_id"]+"/sticker-pack-bans", nil, nil)
		},
	})
	r.add(&Def{
		Name: "stickers.guild_bans.add", Description: "服 ban 贴图包",
		Destructive: true,
		InputSchema: schemaObject(map[string]any{
			"guild_id": propString("服务器 ID"),
			"pack_id":  propString("贴图包 ID"),
			"confirm":  propBool("必须为 true"),
		}, "guild_id", "pack_id", "confirm"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id", "pack_id")
			if err != nil {
				return nil, err
			}
			path := fmt.Sprintf("/guilds/%s/sticker-pack-bans/%s", req["guild_id"], req["pack_id"])
			return c.Gapi(http.MethodPut, path, nil, nil)
		},
	})
	r.add(&Def{
		Name: "stickers.guild_bans.remove", Description: "解除服 ban",
		InputSchema: schemaObject(map[string]any{
			"guild_id": propString("服务器 ID"),
			"pack_id":  propString("贴图包 ID"),
		}, "guild_id", "pack_id"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "guild_id", "pack_id")
			if err != nil {
				return nil, err
			}
			path := fmt.Sprintf("/guilds/%s/sticker-pack-bans/%s", req["guild_id"], req["pack_id"])
			return c.Gapi(http.MethodDelete, path, nil, nil)
		},
	})

	// ---- 社交：隐私 / 好友 / 通知 / 私信 ----
	r.add(&Def{
		Name: "social.privacy.get", Description: "读取隐私设置",
		InputSchema: schemaObject(map[string]any{}),
		CLIHint:     "owl social privacy get",
		run: func(c *api.Client, _ map[string]any) (any, error) {
			return c.Gapi(http.MethodGet, "/users/@me/privacy", nil, nil)
		},
	})
	r.add(&Def{
		Name: "social.privacy.patch", Description: "更新隐私设置",
		InputSchema: schemaObject(map[string]any{
			"friend_request_from":          propString("everyone|mutual_friends|mutual_guilds|nobody"),
			"dm_from":                      propString("everyone|friends|mutual_guilds|nobody"),
			"message_request_filter":       propBool("消息请求过滤"),
			"show_mutual_guilds":           propBool("显示共同服务器"),
			"public_profile_to_non_friends": propBool("非好友可见公开资料"),
		}),
		run: func(c *api.Client, args map[string]any) (any, error) {
			body := bodyFromArgs(args)
			return c.Gapi(http.MethodPatch, "/users/@me/privacy", body, nil)
		},
	})
	r.add(&Def{
		Name: "social.friends.list", Description: "好友/请求/屏蔽关系列表",
		InputSchema: schemaObject(map[string]any{}),
		CLIHint:     "owl social friends list",
		run: func(c *api.Client, _ map[string]any) (any, error) {
			return c.Gapi(http.MethodGet, "/users/@me/relationships", nil, nil)
		},
	})
	r.add(&Def{
		Name: "social.friends.request", Description: "发送好友请求（username 或 user_id）",
		InputSchema: schemaObject(map[string]any{
			"username": propString("用户名"),
			"user_id":  propString("用户 ID"),
		}),
		CLIHint: "owl social friends request --username ...",
		run: func(c *api.Client, args map[string]any) (any, error) {
			body := bodyFromArgs(args)
			if len(body) == 0 {
				return nil, fmt.Errorf("需要 username 或 user_id")
			}
			return c.Gapi(http.MethodPost, "/users/@me/relationships", body, nil)
		},
	})
	r.add(&Def{
		Name: "social.friends.accept", Description: "接受好友请求",
		InputSchema: schemaObject(map[string]any{"user_id": propString("对方用户 ID")}, "user_id"),
		run: func(c *api.Client, args map[string]any) (any, error) {
			req, err := requireStr(args, "user_id")
			if err != nil {
				return nil, err
			}
			return c.Gapi(http.MethodPut, "/users/@me/relationships/"+req["user_id"], map[string]any{"type": "friend"}, nil)
		},
	})
	r.add(&Def{
		Name: "social.friends.block", Description: "屏蔽用户",
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
			return c.Gapi(http.MethodPut, "/users/@me/relationships/"+req["user_id"], map[string]any{"type": "blocked"}, nil)
		},
	})
	r.add(&Def{
		Name: "social.friends.remove", Description: "删除好友/取消请求/解除屏蔽",
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
			return c.Gapi(http.MethodDelete, "/users/@me/relationships/"+req["user_id"], nil, nil)
		},
	})
	r.add(&Def{
		Name: "social.notifications.list", Description: "通知收件箱",
		InputSchema: schemaObject(map[string]any{}),
		CLIHint:     "owl social notifications list",
		run: func(c *api.Client, _ map[string]any) (any, error) {
			return c.Gapi(http.MethodGet, "/users/@me/notifications", nil, nil)
		},
	})
	r.add(&Def{
		Name: "social.notifications.ack", Description: "确认已读通知",
		InputSchema: schemaObject(map[string]any{
			"ids": map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "通知 ID 列表；省略可能表示全部"},
		}),
		run: func(c *api.Client, args map[string]any) (any, error) {
			body := bodyFromArgs(args)
			return c.Gapi(http.MethodPost, "/users/@me/notifications/ack", body, nil)
		},
	})
	r.add(&Def{
		Name: "social.dm.list", Description: "私信/群聊频道列表",
		InputSchema: schemaObject(map[string]any{}),
		CLIHint:     "owl social dm list",
		run: func(c *api.Client, _ map[string]any) (any, error) {
			return c.Gapi(http.MethodGet, "/users/@me/channels", nil, nil)
		},
	})
	r.add(&Def{
		Name: "social.dm.create", Description: "创建 1:1 私信或群聊（recipients 用户 ID 列表）",
		InputSchema: schemaObject(map[string]any{
			"recipients": map[string]any{
				"type": "array", "items": map[string]any{"type": "string"},
				"description": "对方用户 ID；单人=DM，多人=群",
			},
			"name": propString("群聊名称（可选）"),
		}, "recipients"),
		CLIHint: "owl social dm create --users id1,id2",
		run: func(c *api.Client, args map[string]any) (any, error) {
			body := bodyFromArgs(args)
			return c.Gapi(http.MethodPost, "/users/@me/channels", body, nil)
		},
	})
}
