package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

// ---------- restrictions ----------

var restrictionsCmd = &cobra.Command{Use: "restrictions", Short: "多维限制（禁言/禁视等）"}

var restGuild, restUser, restScope, restKind, restChannel, restReason, restExpires, restID string
var restActive string
var restSendText, restViewText, restListen, restSpeak bool

var restrictionsListCmd = &cobra.Command{
	Use: "list", Short: "列出限制",
	Run: func(cmd *cobra.Command, args []string) {
		m := map[string]any{"guild_id": restGuild}
		if restUser != "" {
			m["user_id"] = restUser
		}
		if restChannel != "" {
			m["channel_id"] = restChannel
		}
		if restScope != "" {
			m["scope"] = restScope
		}
		if restActive != "" {
			m["active"] = restActive
		}
		runTool("restrictions.list", m, flagYes)
	},
}

var restrictionsCreateCmd = &cobra.Command{
	Use: "create", Short: "创建限制（危险）",
	Run: func(cmd *cobra.Command, args []string) {
		deny := map[string]any{}
		if restSendText {
			deny["send_text"] = true
		}
		if restViewText {
			deny["view_text"] = true
		}
		if restListen {
			deny["listen_voice"] = true
		}
		if restSpeak {
			deny["speak_voice"] = true
		}
		if len(deny) == 0 {
			deny["send_text"] = true
		}
		m := map[string]any{
			"guild_id": restGuild, "target_user_id": restUser,
			"scope": restScope, "kind": restKind, "deny": deny, "confirm": true,
		}
		if restChannel != "" {
			m["channel_id"] = restChannel
		}
		if restReason != "" {
			m["reason"] = restReason
		}
		if restExpires != "" {
			m["expires_at"] = restExpires
		}
		runTool("restrictions.create", m, flagYes)
	},
}

var restrictionsLiftCmd = &cobra.Command{
	Use: "lift", Short: "解除限制（危险）",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("restrictions.lift", map[string]any{
			"guild_id": restGuild, "restriction_id": restID, "confirm": true,
		}, flagYes)
	},
}

// ---------- audit ----------

var auditCmd = &cobra.Command{Use: "audit", Short: "审计日志"}

var auditGuild, auditAction string
var auditLimit int

var auditListCmd = &cobra.Command{
	Use: "list", Short: "服务器审计",
	Run: func(cmd *cobra.Command, args []string) {
		m := map[string]any{"guild_id": auditGuild}
		if auditAction != "" {
			m["action"] = auditAction
		}
		if auditLimit > 0 {
			m["limit"] = auditLimit
		}
		runTool("audit.list", m, flagYes)
	},
}

// ---------- voice ----------

var voiceCmd = &cobra.Command{Use: "voice", Short: "语音治理与节点池"}

var voiceGuild, voiceChannel, voiceUser string
var voiceMute, voiceDeaf bool
var voiceNodes string
var voiceFallback bool

var voiceStatesCmd = &cobra.Command{
	Use: "states", Short: "频道语音状态",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("voice.states", map[string]any{
			"guild_id": voiceGuild, "channel_id": voiceChannel,
		}, flagYes)
	},
}

var voiceDisconnectCmd = &cobra.Command{
	Use: "disconnect", Short: "踢出语音（危险）",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("voice.disconnect", map[string]any{
			"guild_id": voiceGuild, "user_id": voiceUser, "confirm": true,
		}, flagYes)
	},
}

var voiceMoveCmd = &cobra.Command{
	Use: "move", Short: "移动语音成员（危险）",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("voice.move", map[string]any{
			"guild_id": voiceGuild, "user_id": voiceUser, "channel_id": voiceChannel, "confirm": true,
		}, flagYes)
	},
}

var voiceMuteCmd = &cobra.Command{
	Use: "mute", Short: "服务器静音/耳聋（危险）",
	Run: func(cmd *cobra.Command, args []string) {
		m := map[string]any{"guild_id": voiceGuild, "user_id": voiceUser, "confirm": true}
		if cmd.Flags().Changed("mute") {
			m["server_mute"] = voiceMute
		}
		if cmd.Flags().Changed("deafen") {
			m["server_deaf"] = voiceDeaf
		}
		runTool("voice.server_mute", m, flagYes)
	},
}

var voiceNodePoolCmd = &cobra.Command{Use: "node-pool", Short: "服务器节点池"}

var voiceNodePoolGetCmd = &cobra.Command{
	Use: "get", Short: "查看节点池",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("voice.node_pool.get", map[string]any{"guild_id": voiceGuild}, flagYes)
	},
}

var voiceNodePoolSetCmd = &cobra.Command{
	Use: "set", Short: "设置生效节点",
	Run: func(cmd *cobra.Command, args []string) {
		var ids []string
		for _, p := range strings.Split(voiceNodes, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				ids = append(ids, p)
			}
		}
		m := map[string]any{"guild_id": voiceGuild, "node_ids": ids}
		if cmd.Flags().Changed("fallback") {
			m["fallback_to_default"] = voiceFallback
		}
		runTool("voice.node_pool.set", m, flagYes)
	},
}

// ---------- platform ----------

var platformCmd = &cobra.Command{Use: "platform", Short: "平台管理（需 --platform 登录 + system_admin）"}

var platformUsersCmd = &cobra.Command{Use: "users", Short: "平台用户"}
var platQ, platUserID, platPassword string
var platAdmin bool

var platformUsersListCmd = &cobra.Command{
	Use: "list", Short: "用户目录",
	Run: func(cmd *cobra.Command, args []string) {
		m := map[string]any{}
		if platQ != "" {
			m["q"] = platQ
		}
		runTool("platform.users.list", m, flagYes)
	},
}

var platformUsersDisableCmd = &cobra.Command{
	Use: "disable", Short: "禁用用户（危险）",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("platform.users.disable", map[string]any{"user_id": platUserID, "confirm": true}, flagYes)
	},
}

var platformUsersEnableCmd = &cobra.Command{
	Use: "enable", Short: "启用用户",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("platform.users.enable", map[string]any{"user_id": platUserID}, flagYes)
	},
}

var platformRegCmd = &cobra.Command{Use: "registration", Short: "注册开关"}
var platSignup bool

var platformRegGetCmd = &cobra.Command{
	Use: "get", Short: "读取注册开关",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("platform.registration.get", map[string]any{}, flagYes)
	},
}

var platformRegSetCmd = &cobra.Command{
	Use: "set", Short: "设置注册开关（危险）",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("platform.registration.set", map[string]any{
			"signup_enabled": platSignup, "confirm": true,
		}, flagYes)
	},
}

var platformSFUCmd = &cobra.Command{Use: "sfu", Short: "SFU 节点"}
var platformSFUNodesCmd = &cobra.Command{
	Use: "nodes", Short: "列出节点",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("platform.sfu.nodes", map[string]any{}, flagYes)
	},
}
var platformSFUTopoCmd = &cobra.Command{
	Use: "topology", Short: "拓扑",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("platform.sfu.topology", map[string]any{}, flagYes)
	},
}

// ---------- roles update 补齐 ----------

var rolesUpdateCmd = &cobra.Command{
	Use: "update", Short: "更新角色",
	Run: func(cmd *cobra.Command, args []string) {
		m := map[string]any{
			"guild_id": rolesGuild, "role_id": roleID, "name": roleName,
			"permissions": rolePerms, "position": rolePos,
		}
		if roleColor != "" {
			m["color"] = roleColor
		}
		runTool("roles.update", m, flagYes)
	},
}

func initExtraCommands() {
	// restrictions
	restrictionsListCmd.Flags().StringVar(&restGuild, "guild", "", "服务器 ID")
	restrictionsListCmd.Flags().StringVar(&restUser, "user", "", "用户过滤")
	restrictionsListCmd.Flags().StringVar(&restChannel, "channel", "", "频道过滤")
	restrictionsListCmd.Flags().StringVar(&restScope, "scope", "", "范围")
	restrictionsListCmd.Flags().StringVar(&restActive, "active", "", "true|false")
	_ = restrictionsListCmd.MarkFlagRequired("guild")
	restrictionsCreateCmd.Flags().StringVar(&restGuild, "guild", "", "服务器 ID")
	restrictionsCreateCmd.Flags().StringVar(&restUser, "user", "", "目标用户")
	restrictionsCreateCmd.Flags().StringVar(&restScope, "scope", "GUILD_ALL_TEXT", "范围")
	restrictionsCreateCmd.Flags().StringVar(&restKind, "kind", "SANCTION", "SANCTION|CHANNEL_BAN")
	restrictionsCreateCmd.Flags().StringVar(&restChannel, "channel", "", "频道 ID")
	restrictionsCreateCmd.Flags().StringVar(&restReason, "reason", "", "原因")
	restrictionsCreateCmd.Flags().StringVar(&restExpires, "expires", "", "RFC3339")
	restrictionsCreateCmd.Flags().BoolVar(&restSendText, "deny-send", true, "禁发文字")
	restrictionsCreateCmd.Flags().BoolVar(&restViewText, "deny-view", false, "禁看文字")
	restrictionsCreateCmd.Flags().BoolVar(&restListen, "deny-listen", false, "禁听语音")
	restrictionsCreateCmd.Flags().BoolVar(&restSpeak, "deny-speak", false, "禁说语音")
	_ = restrictionsCreateCmd.MarkFlagRequired("guild")
	_ = restrictionsCreateCmd.MarkFlagRequired("user")
	restrictionsLiftCmd.Flags().StringVar(&restGuild, "guild", "", "服务器 ID")
	restrictionsLiftCmd.Flags().StringVar(&restID, "id", "", "限制 ID")
	_ = restrictionsLiftCmd.MarkFlagRequired("guild")
	_ = restrictionsLiftCmd.MarkFlagRequired("id")
	restrictionsCmd.AddCommand(restrictionsListCmd, restrictionsCreateCmd, restrictionsLiftCmd)

	// audit
	auditListCmd.Flags().StringVar(&auditGuild, "guild", "", "服务器 ID")
	auditListCmd.Flags().StringVar(&auditAction, "action", "", "action 前缀")
	auditListCmd.Flags().IntVar(&auditLimit, "limit", 50, "条数")
	_ = auditListCmd.MarkFlagRequired("guild")
	auditCmd.AddCommand(auditListCmd)

	// voice
	voiceStatesCmd.Flags().StringVar(&voiceGuild, "guild", "", "服务器 ID")
	voiceStatesCmd.Flags().StringVar(&voiceChannel, "channel", "", "频道 ID")
	_ = voiceStatesCmd.MarkFlagRequired("guild")
	_ = voiceStatesCmd.MarkFlagRequired("channel")
	voiceDisconnectCmd.Flags().StringVar(&voiceGuild, "guild", "", "服务器 ID")
	voiceDisconnectCmd.Flags().StringVar(&voiceUser, "user", "", "用户 ID")
	_ = voiceDisconnectCmd.MarkFlagRequired("guild")
	_ = voiceDisconnectCmd.MarkFlagRequired("user")
	voiceMoveCmd.Flags().StringVar(&voiceGuild, "guild", "", "服务器 ID")
	voiceMoveCmd.Flags().StringVar(&voiceUser, "user", "", "用户 ID")
	voiceMoveCmd.Flags().StringVar(&voiceChannel, "channel", "", "目标频道")
	_ = voiceMoveCmd.MarkFlagRequired("guild")
	_ = voiceMoveCmd.MarkFlagRequired("user")
	_ = voiceMoveCmd.MarkFlagRequired("channel")
	voiceMuteCmd.Flags().StringVar(&voiceGuild, "guild", "", "服务器 ID")
	voiceMuteCmd.Flags().StringVar(&voiceUser, "user", "", "用户 ID")
	voiceMuteCmd.Flags().BoolVar(&voiceMute, "mute", false, "服务器静音")
	voiceMuteCmd.Flags().BoolVar(&voiceDeaf, "deafen", false, "服务器耳聋")
	_ = voiceMuteCmd.MarkFlagRequired("guild")
	_ = voiceMuteCmd.MarkFlagRequired("user")
	voiceNodePoolGetCmd.Flags().StringVar(&voiceGuild, "guild", "", "服务器 ID")
	_ = voiceNodePoolGetCmd.MarkFlagRequired("guild")
	voiceNodePoolSetCmd.Flags().StringVar(&voiceGuild, "guild", "", "服务器 ID")
	voiceNodePoolSetCmd.Flags().StringVar(&voiceNodes, "nodes", "", "逗号分隔节点 ID")
	voiceNodePoolSetCmd.Flags().BoolVar(&voiceFallback, "fallback", false, "回退默认池")
	_ = voiceNodePoolSetCmd.MarkFlagRequired("guild")
	_ = voiceNodePoolSetCmd.MarkFlagRequired("nodes")
	voiceNodePoolCmd.AddCommand(voiceNodePoolGetCmd, voiceNodePoolSetCmd)
	voiceCmd.AddCommand(voiceStatesCmd, voiceDisconnectCmd, voiceMoveCmd, voiceMuteCmd, voiceNodePoolCmd)

	// platform
	platformUsersListCmd.Flags().StringVar(&platQ, "q", "", "搜索")
	platformUsersDisableCmd.Flags().StringVar(&platUserID, "user", "", "用户 ID")
	_ = platformUsersDisableCmd.MarkFlagRequired("user")
	platformUsersEnableCmd.Flags().StringVar(&platUserID, "user", "", "用户 ID")
	_ = platformUsersEnableCmd.MarkFlagRequired("user")
	platformUsersCmd.AddCommand(platformUsersListCmd, platformUsersDisableCmd, platformUsersEnableCmd)
	platformRegSetCmd.Flags().BoolVar(&platSignup, "enabled", true, "开放注册")
	platformRegCmd.AddCommand(platformRegGetCmd, platformRegSetCmd)
	platformSFUCmd.AddCommand(platformSFUNodesCmd, platformSFUTopoCmd)
	platformCmd.AddCommand(platformUsersCmd, platformRegCmd, platformSFUCmd)

	// roles update
	rolesUpdateCmd.Flags().StringVar(&rolesGuild, "guild", "", "服务器 ID")
	rolesUpdateCmd.Flags().StringVar(&roleID, "role", "", "角色 ID")
	rolesUpdateCmd.Flags().StringVar(&roleName, "name", "", "名称")
	rolesUpdateCmd.Flags().Int64Var(&rolePerms, "permissions", 0, "权限位")
	rolesUpdateCmd.Flags().IntVar(&rolePos, "position", 1, "位置")
	rolesUpdateCmd.Flags().StringVar(&roleColor, "color", "", "颜色")
	_ = rolesUpdateCmd.MarkFlagRequired("guild")
	_ = rolesUpdateCmd.MarkFlagRequired("role")
	_ = rolesUpdateCmd.MarkFlagRequired("name")
	rolesCmd.AddCommand(rolesUpdateCmd)
}
