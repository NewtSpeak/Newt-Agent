package cmd

import (
	"encoding/json"
	"strconv"

	"github.com/spf13/cobra"
)

// 共享 --yes 标志（危险操作）
var flagYes bool

func bindYes(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&flagYes, "yes", false, "跳过危险操作确认（等同 confirm=true）")
}

// ---------- guilds 扩展 ----------

var guildsGetCmd = &cobra.Command{
	Use:   "get <guild_id>",
	Short: "服务器详情",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runTool("guilds.get", map[string]any{"guild_id": args[0]}, flagYes)
	},
}

var (
	guildCreateName string
)

var guildsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建服务器",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("guilds.create", map[string]any{"name": guildCreateName}, flagYes)
	},
}

var (
	guildUpdateID, guildUpdateName, guildUpdateDesc, guildUpdateDefaultCh string
)

var guildsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "更新服务器",
	Run: func(cmd *cobra.Command, args []string) {
		m := map[string]any{"guild_id": guildUpdateID}
		if guildUpdateName != "" {
			m["name"] = guildUpdateName
		}
		if guildUpdateDesc != "" {
			m["description"] = guildUpdateDesc
		}
		if cmd.Flags().Changed("default-channel") {
			m["default_channel_id"] = guildUpdateDefaultCh
		}
		runTool("guilds.update", m, flagYes)
	},
}

var (
	guildDeleteID, guildDeleteConfirmName string
)

var guildsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "删除服务器（危险）",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("guilds.delete", map[string]any{
			"guild_id": guildDeleteID, "confirm_name": guildDeleteConfirmName, "confirm": true,
		}, flagYes)
	},
}

var guildsPermsCmd = &cobra.Command{
	Use:   "permissions <guild_id>",
	Short: "本人在服务器的权限位",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runTool("guilds.permissions", map[string]any{"guild_id": args[0]}, flagYes)
	},
}

// ---------- channels ----------

var channelsCmd = &cobra.Command{
	Use:   "channels",
	Short: "频道管理",
}

var channelsGuild string

var channelsListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出频道",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("channels.list", map[string]any{"guild_id": channelsGuild}, flagYes)
	},
}

var (
	chCreateName, chCreateType, chCreateTopic, chCreateParent, chCreatePassword string
	chCreatePrivate                                                            bool
)

var channelsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建频道",
	Run: func(cmd *cobra.Command, args []string) {
		m := map[string]any{
			"guild_id": channelsGuild,
			"name":     chCreateName,
			"type":     chCreateType,
		}
		if chCreateTopic != "" {
			m["topic"] = chCreateTopic
		}
		if chCreateParent != "" {
			m["parent_id"] = chCreateParent
		}
		if chCreatePassword != "" {
			m["password"] = chCreatePassword
		}
		if chCreatePrivate {
			m["private"] = true
		}
		runTool("channels.create", m, flagYes)
	},
}

var (
	chUpdateID, chUpdateName, chUpdateTopic, chUpdateParent string
)

var channelsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "更新频道",
	Run: func(cmd *cobra.Command, args []string) {
		m := map[string]any{"channel_id": chUpdateID}
		if chUpdateName != "" {
			m["name"] = chUpdateName
		}
		if chUpdateTopic != "" {
			m["topic"] = chUpdateTopic
		}
		if cmd.Flags().Changed("parent") {
			m["parent_id"] = chUpdateParent
		}
		runTool("channels.update", m, flagYes)
	},
}

var chDeleteID string

var channelsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "删除频道（危险）",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("channels.delete", map[string]any{"channel_id": chDeleteID, "confirm": true}, flagYes)
	},
}

var (
	owChannel, owTarget, owType string
	owAllow, owDeny             int64
)

var channelsOverwritesCmd = &cobra.Command{
	Use:   "overwrites",
	Short: "频道权限覆盖",
}

var owGuild, owListChannel string

var overwritesListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出覆盖",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("channels.overwrites.list", map[string]any{
			"guild_id": owGuild, "channel_id": owListChannel,
		}, flagYes)
	},
}

var overwritesSetCmd = &cobra.Command{
	Use:   "set",
	Short: "设置覆盖",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("channels.overwrites.upsert", map[string]any{
			"channel_id": owChannel, "target_id": owTarget, "type": owType,
			"allow": owAllow, "deny": owDeny,
		}, flagYes)
	},
}

var overwritesDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "删除覆盖（危险）",
	Run: func(cmd *cobra.Command, args []string) {
		m := map[string]any{"channel_id": owChannel, "target_id": owTarget, "confirm": true}
		if owType != "" {
			m["type"] = owType
		}
		runTool("channels.overwrites.delete", m, flagYes)
	},
}

// ---------- roles ----------

var rolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "角色管理",
}

var rolesGuild string

var rolesListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出角色",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("roles.list", map[string]any{"guild_id": rolesGuild}, flagYes)
	},
}

var (
	roleName, roleColor string
	rolePerms           int64
	rolePos             int
	roleHoist, roleMentionable bool
)

var rolesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建角色",
	Run: func(cmd *cobra.Command, args []string) {
		m := map[string]any{
			"guild_id": rolesGuild, "name": roleName,
			"permissions": rolePerms, "position": rolePos,
		}
		if roleColor != "" {
			m["color"] = roleColor
		}
		if roleHoist {
			m["hoist"] = true
		}
		if roleMentionable {
			m["mentionable"] = true
		}
		runTool("roles.create", m, flagYes)
	},
}

var roleID string

var rolesDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "删除角色（危险）",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("roles.delete", map[string]any{
			"guild_id": rolesGuild, "role_id": roleID, "confirm": true,
		}, flagYes)
	},
}

var roleMemberID string

var rolesAssignCmd = &cobra.Command{
	Use:   "assign",
	Short: "赋角",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("roles.assign", map[string]any{
			"guild_id": rolesGuild, "member_id": roleMemberID, "role_id": roleID,
		}, flagYes)
	},
}

var rolesRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "摘角（危险）",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("roles.remove", map[string]any{
			"guild_id": rolesGuild, "member_id": roleMemberID, "role_id": roleID, "confirm": true,
		}, flagYes)
	},
}

// ---------- members ----------

var membersCmd = &cobra.Command{
	Use:   "members",
	Short: "成员治理",
}

var membersGuild string

var membersListCmd = &cobra.Command{
	Use:   "list",
	Short: "成员列表",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("members.list", map[string]any{"guild_id": membersGuild}, flagYes)
	},
}

var memberID, memberNick string

var membersNickCmd = &cobra.Command{
	Use:   "nick",
	Short: "改昵称",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("members.nick", map[string]any{
			"guild_id": membersGuild, "member_id": memberID, "nick": memberNick,
		}, flagYes)
	},
}

var membersKickCmd = &cobra.Command{
	Use:   "kick",
	Short: "踢出（危险）",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("members.kick", map[string]any{
			"guild_id": membersGuild, "member_id": memberID, "confirm": true,
		}, flagYes)
	},
}

var banUserID, banReason string

var membersBanCmd = &cobra.Command{
	Use:   "ban",
	Short: "封禁（危险）",
	Run: func(cmd *cobra.Command, args []string) {
		m := map[string]any{"guild_id": membersGuild, "user_id": banUserID, "confirm": true}
		if banReason != "" {
			m["reason"] = banReason
		}
		runTool("members.ban", m, flagYes)
	},
}

var membersUnbanCmd = &cobra.Command{
	Use:   "unban",
	Short: "解封",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("members.unban", map[string]any{"guild_id": membersGuild, "user_id": banUserID}, flagYes)
	},
}

var membersBansCmd = &cobra.Command{
	Use:   "bans",
	Short: "封禁列表",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("members.bans", map[string]any{"guild_id": membersGuild}, flagYes)
	},
}

// ---------- invites ----------

var invitesCmd = &cobra.Command{
	Use:   "invites",
	Short: "邀请",
}

var inviteGuild string
var inviteTTL, inviteMaxUses int

var invitesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建邀请",
	Run: func(cmd *cobra.Command, args []string) {
		m := map[string]any{"guild_id": inviteGuild}
		if inviteTTL >= 60 {
			m["ttl_seconds"] = inviteTTL
		}
		if inviteMaxUses >= 1 {
			m["max_uses"] = inviteMaxUses
		}
		runTool("invites.create", m, flagYes)
	},
}

var inviteCode string

var invitesGetCmd = &cobra.Command{
	Use:   "get <code>",
	Short: "预览邀请",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runTool("invites.get", map[string]any{"code": args[0]}, flagYes)
	},
}

var invitesJoinCmd = &cobra.Command{
	Use:   "join <code>",
	Short: "加入邀请",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runTool("invites.join", map[string]any{"code": args[0]}, flagYes)
	},
}

// ---------- messages ----------

var messagesCmd = &cobra.Command{
	Use:   "messages",
	Short: "消息",
}

var msgChannel, msgContent, msgID, msgSearchQ string
var msgLimit int

var messagesListCmd = &cobra.Command{
	Use:   "list",
	Short: "历史消息",
	Run: func(cmd *cobra.Command, args []string) {
		m := map[string]any{"channel_id": msgChannel}
		if msgLimit > 0 {
			m["limit"] = msgLimit
		}
		runTool("messages.list", m, flagYes)
	},
}

var messagesSendCmd = &cobra.Command{
	Use:   "send",
	Short: "发消息",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("messages.send", map[string]any{"channel_id": msgChannel, "content": msgContent}, flagYes)
	},
}

var messagesGetCmd = &cobra.Command{
	Use:   "get",
	Short: "单条消息",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("messages.get", map[string]any{"channel_id": msgChannel, "message_id": msgID}, flagYes)
	},
}

var messagesDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "删消息（危险）",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("messages.delete", map[string]any{
			"channel_id": msgChannel, "message_id": msgID, "confirm": true,
		}, flagYes)
	},
}

var messagesSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "搜索消息",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("messages.search", map[string]any{"q": msgSearchQ}, flagYes)
	},
}

func initAdminCommands() {
	// guilds
	guildsCreateCmd.Flags().StringVar(&guildCreateName, "name", "", "服务器名称")
	_ = guildsCreateCmd.MarkFlagRequired("name")
	guildsUpdateCmd.Flags().StringVar(&guildUpdateID, "guild", "", "服务器 ID")
	guildsUpdateCmd.Flags().StringVar(&guildUpdateName, "name", "", "新名称")
	guildsUpdateCmd.Flags().StringVar(&guildUpdateDesc, "description", "", "简介")
	guildsUpdateCmd.Flags().StringVar(&guildUpdateDefaultCh, "default-channel", "", "默认频道")
	_ = guildsUpdateCmd.MarkFlagRequired("guild")
	guildsDeleteCmd.Flags().StringVar(&guildDeleteID, "guild", "", "服务器 ID")
	guildsDeleteCmd.Flags().StringVar(&guildDeleteConfirmName, "confirm-name", "", "确认名称")
	_ = guildsDeleteCmd.MarkFlagRequired("guild")
	_ = guildsDeleteCmd.MarkFlagRequired("confirm-name")
	guildsCmd.AddCommand(guildsGetCmd, guildsCreateCmd, guildsUpdateCmd, guildsDeleteCmd, guildsPermsCmd)

	// channels
	channelsListCmd.Flags().StringVar(&channelsGuild, "guild", "", "服务器 ID")
	_ = channelsListCmd.MarkFlagRequired("guild")
	channelsCreateCmd.Flags().StringVar(&channelsGuild, "guild", "", "服务器 ID")
	channelsCreateCmd.Flags().StringVar(&chCreateName, "name", "", "频道名")
	channelsCreateCmd.Flags().StringVar(&chCreateType, "type", "TEXT", "TEXT|VOICE|CATEGORY|STAGE")
	channelsCreateCmd.Flags().StringVar(&chCreateTopic, "topic", "", "主题")
	channelsCreateCmd.Flags().StringVar(&chCreateParent, "parent", "", "父分类")
	channelsCreateCmd.Flags().StringVar(&chCreatePassword, "password", "", "访问密码")
	channelsCreateCmd.Flags().BoolVar(&chCreatePrivate, "private", false, "私密")
	_ = channelsCreateCmd.MarkFlagRequired("guild")
	_ = channelsCreateCmd.MarkFlagRequired("name")
	channelsUpdateCmd.Flags().StringVar(&chUpdateID, "channel", "", "频道 ID")
	channelsUpdateCmd.Flags().StringVar(&chUpdateName, "name", "", "名称")
	channelsUpdateCmd.Flags().StringVar(&chUpdateTopic, "topic", "", "主题")
	channelsUpdateCmd.Flags().StringVar(&chUpdateParent, "parent", "", "父分类")
	_ = channelsUpdateCmd.MarkFlagRequired("channel")
	channelsDeleteCmd.Flags().StringVar(&chDeleteID, "channel", "", "频道 ID")
	_ = channelsDeleteCmd.MarkFlagRequired("channel")
	overwritesListCmd.Flags().StringVar(&owGuild, "guild", "", "服务器 ID")
	overwritesListCmd.Flags().StringVar(&owListChannel, "channel", "", "频道 ID")
	_ = overwritesListCmd.MarkFlagRequired("guild")
	_ = overwritesListCmd.MarkFlagRequired("channel")
	overwritesSetCmd.Flags().StringVar(&owChannel, "channel", "", "频道 ID")
	overwritesSetCmd.Flags().StringVar(&owTarget, "target", "", "角色/成员 ID")
	overwritesSetCmd.Flags().StringVar(&owType, "type", "ROLE", "ROLE|MEMBER")
	overwritesSetCmd.Flags().Int64Var(&owAllow, "allow", 0, "允许位")
	overwritesSetCmd.Flags().Int64Var(&owDeny, "deny", 0, "拒绝位")
	_ = overwritesSetCmd.MarkFlagRequired("channel")
	_ = overwritesSetCmd.MarkFlagRequired("target")
	overwritesDeleteCmd.Flags().StringVar(&owChannel, "channel", "", "频道 ID")
	overwritesDeleteCmd.Flags().StringVar(&owTarget, "target", "", "目标 ID")
	overwritesDeleteCmd.Flags().StringVar(&owType, "type", "", "ROLE|MEMBER")
	_ = overwritesDeleteCmd.MarkFlagRequired("channel")
	_ = overwritesDeleteCmd.MarkFlagRequired("target")
	channelsOverwritesCmd.AddCommand(overwritesListCmd, overwritesSetCmd, overwritesDeleteCmd)
	channelsCmd.AddCommand(channelsListCmd, channelsCreateCmd, channelsUpdateCmd, channelsDeleteCmd, channelsOverwritesCmd)

	// roles
	rolesListCmd.Flags().StringVar(&rolesGuild, "guild", "", "服务器 ID")
	_ = rolesListCmd.MarkFlagRequired("guild")
	rolesCreateCmd.Flags().StringVar(&rolesGuild, "guild", "", "服务器 ID")
	rolesCreateCmd.Flags().StringVar(&roleName, "name", "", "角色名")
	rolesCreateCmd.Flags().Int64Var(&rolePerms, "permissions", 0, "权限位")
	rolesCreateCmd.Flags().IntVar(&rolePos, "position", 1, "位置")
	rolesCreateCmd.Flags().StringVar(&roleColor, "color", "", "颜色")
	rolesCreateCmd.Flags().BoolVar(&roleHoist, "hoist", false, "单独显示")
	rolesCreateCmd.Flags().BoolVar(&roleMentionable, "mentionable", false, "可提及")
	_ = rolesCreateCmd.MarkFlagRequired("guild")
	_ = rolesCreateCmd.MarkFlagRequired("name")
	rolesDeleteCmd.Flags().StringVar(&rolesGuild, "guild", "", "服务器 ID")
	rolesDeleteCmd.Flags().StringVar(&roleID, "role", "", "角色 ID")
	_ = rolesDeleteCmd.MarkFlagRequired("guild")
	_ = rolesDeleteCmd.MarkFlagRequired("role")
	rolesAssignCmd.Flags().StringVar(&rolesGuild, "guild", "", "服务器 ID")
	rolesAssignCmd.Flags().StringVar(&roleMemberID, "member", "", "成员 ID")
	rolesAssignCmd.Flags().StringVar(&roleID, "role", "", "角色 ID")
	_ = rolesAssignCmd.MarkFlagRequired("guild")
	_ = rolesAssignCmd.MarkFlagRequired("member")
	_ = rolesAssignCmd.MarkFlagRequired("role")
	rolesRemoveCmd.Flags().StringVar(&rolesGuild, "guild", "", "服务器 ID")
	rolesRemoveCmd.Flags().StringVar(&roleMemberID, "member", "", "成员 ID")
	rolesRemoveCmd.Flags().StringVar(&roleID, "role", "", "角色 ID")
	_ = rolesRemoveCmd.MarkFlagRequired("guild")
	_ = rolesRemoveCmd.MarkFlagRequired("member")
	_ = rolesRemoveCmd.MarkFlagRequired("role")
	rolesCmd.AddCommand(rolesListCmd, rolesCreateCmd, rolesDeleteCmd, rolesAssignCmd, rolesRemoveCmd)

	// members
	membersListCmd.Flags().StringVar(&membersGuild, "guild", "", "服务器 ID")
	_ = membersListCmd.MarkFlagRequired("guild")
	membersNickCmd.Flags().StringVar(&membersGuild, "guild", "", "服务器 ID")
	membersNickCmd.Flags().StringVar(&memberID, "member", "", "成员 ID")
	membersNickCmd.Flags().StringVar(&memberNick, "nick", "", "昵称")
	_ = membersNickCmd.MarkFlagRequired("guild")
	_ = membersNickCmd.MarkFlagRequired("member")
	membersKickCmd.Flags().StringVar(&membersGuild, "guild", "", "服务器 ID")
	membersKickCmd.Flags().StringVar(&memberID, "member", "", "成员 ID")
	_ = membersKickCmd.MarkFlagRequired("guild")
	_ = membersKickCmd.MarkFlagRequired("member")
	membersBanCmd.Flags().StringVar(&membersGuild, "guild", "", "服务器 ID")
	membersBanCmd.Flags().StringVar(&banUserID, "user", "", "用户 ID")
	membersBanCmd.Flags().StringVar(&banReason, "reason", "", "原因")
	_ = membersBanCmd.MarkFlagRequired("guild")
	_ = membersBanCmd.MarkFlagRequired("user")
	membersUnbanCmd.Flags().StringVar(&membersGuild, "guild", "", "服务器 ID")
	membersUnbanCmd.Flags().StringVar(&banUserID, "user", "", "用户 ID")
	_ = membersUnbanCmd.MarkFlagRequired("guild")
	_ = membersUnbanCmd.MarkFlagRequired("user")
	membersBansCmd.Flags().StringVar(&membersGuild, "guild", "", "服务器 ID")
	_ = membersBansCmd.MarkFlagRequired("guild")
	membersCmd.AddCommand(membersListCmd, membersNickCmd, membersKickCmd, membersBanCmd, membersUnbanCmd, membersBansCmd)

	// invites
	invitesCreateCmd.Flags().StringVar(&inviteGuild, "guild", "", "服务器 ID")
	invitesCreateCmd.Flags().IntVar(&inviteTTL, "ttl", 0, "有效秒数 ≥60")
	invitesCreateCmd.Flags().IntVar(&inviteMaxUses, "max-uses", 0, "最大次数")
	_ = invitesCreateCmd.MarkFlagRequired("guild")
	invitesCmd.AddCommand(invitesCreateCmd, invitesGetCmd, invitesJoinCmd)

	// messages
	messagesListCmd.Flags().StringVar(&msgChannel, "channel", "", "频道 ID")
	messagesListCmd.Flags().IntVar(&msgLimit, "limit", 0, "条数")
	_ = messagesListCmd.MarkFlagRequired("channel")
	messagesSendCmd.Flags().StringVar(&msgChannel, "channel", "", "频道 ID")
	messagesSendCmd.Flags().StringVar(&msgContent, "content", "", "正文")
	_ = messagesSendCmd.MarkFlagRequired("channel")
	_ = messagesSendCmd.MarkFlagRequired("content")
	messagesGetCmd.Flags().StringVar(&msgChannel, "channel", "", "频道 ID")
	messagesGetCmd.Flags().StringVar(&msgID, "message", "", "消息 ID")
	_ = messagesGetCmd.MarkFlagRequired("channel")
	_ = messagesGetCmd.MarkFlagRequired("message")
	messagesDeleteCmd.Flags().StringVar(&msgChannel, "channel", "", "频道 ID")
	messagesDeleteCmd.Flags().StringVar(&msgID, "message", "", "消息 ID")
	_ = messagesDeleteCmd.MarkFlagRequired("channel")
	_ = messagesDeleteCmd.MarkFlagRequired("message")
	messagesSearchCmd.Flags().StringVar(&msgSearchQ, "q", "", "关键词")
	_ = messagesSearchCmd.MarkFlagRequired("q")
	messagesCmd.AddCommand(messagesListCmd, messagesSendCmd, messagesGetCmd, messagesDeleteCmd, messagesSearchCmd)

	_ = json.Marshal
	_ = strconv.Itoa
}
