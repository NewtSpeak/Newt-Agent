package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

// ---------- stickers ----------

var stickersCmd = &cobra.Command{Use: "stickers", Short: "贴图包与贴图库"}

var stickersPacksCmd = &cobra.Command{Use: "packs", Short: "自建贴图包"}
var stickersLibraryCmd = &cobra.Command{Use: "library", Short: "贴图库"}
var stickerPackName, stickerPackID, stickerGuild string

var stickersPacksListCmd = &cobra.Command{
	Use: "list", Short: "我的贴图包",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("stickers.packs.list", map[string]any{}, flagYes)
	},
}
var stickersPacksCreateCmd = &cobra.Command{
	Use: "create", Short: "创建贴图包",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("stickers.packs.create", map[string]any{"name": stickerPackName}, flagYes)
	},
}
var stickersLibraryListCmd = &cobra.Command{
	Use: "list", Short: "已安装库",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("stickers.library.list", map[string]any{}, flagYes)
	},
}
var stickersAvailableCmd = &cobra.Command{
	Use: "available", Short: "可用贴图集合",
	Run: func(cmd *cobra.Command, args []string) {
		m := map[string]any{}
		if stickerGuild != "" {
			m["guild_id"] = stickerGuild
		}
		runTool("stickers.available", m, flagYes)
	},
}

// ---------- social ----------

var socialCmd = &cobra.Command{Use: "social", Short: "好友 / 隐私 / 通知 / 私信"}

var socialPrivacyCmd = &cobra.Command{Use: "privacy", Short: "隐私设置"}
var socialFriendsCmd = &cobra.Command{Use: "friends", Short: "好友关系"}
var socialNotifCmd = &cobra.Command{Use: "notifications", Short: "通知"}
var socialDMCmd = &cobra.Command{Use: "dm", Short: "私信"}

var socialUsername, socialUserID, socialDMUsers string

var socialPrivacyGetCmd = &cobra.Command{
	Use: "get", Short: "读取隐私",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("social.privacy.get", map[string]any{}, flagYes)
	},
}
var socialFriendsListCmd = &cobra.Command{
	Use: "list", Short: "关系列表",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("social.friends.list", map[string]any{}, flagYes)
	},
}
var socialFriendsRequestCmd = &cobra.Command{
	Use: "request", Short: "发好友请求",
	Run: func(cmd *cobra.Command, args []string) {
		m := map[string]any{}
		if socialUsername != "" {
			m["username"] = socialUsername
		}
		if socialUserID != "" {
			m["user_id"] = socialUserID
		}
		runTool("social.friends.request", m, flagYes)
	},
}
var socialFriendsAcceptCmd = &cobra.Command{
	Use: "accept", Short: "接受请求",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("social.friends.accept", map[string]any{"user_id": socialUserID}, flagYes)
	},
}
var socialFriendsBlockCmd = &cobra.Command{
	Use: "block", Short: "屏蔽（危险）",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("social.friends.block", map[string]any{"user_id": socialUserID, "confirm": true}, flagYes)
	},
}
var socialFriendsRemoveCmd = &cobra.Command{
	Use: "remove", Short: "删除关系（危险）",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("social.friends.remove", map[string]any{"user_id": socialUserID, "confirm": true}, flagYes)
	},
}
var socialNotifListCmd = &cobra.Command{
	Use: "list", Short: "通知列表",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("social.notifications.list", map[string]any{}, flagYes)
	},
}
var socialDMListCmd = &cobra.Command{
	Use: "list", Short: "私信频道",
	Run: func(cmd *cobra.Command, args []string) {
		runTool("social.dm.list", map[string]any{}, flagYes)
	},
}
var socialDMCreateCmd = &cobra.Command{
	Use: "create", Short: "创建私信/群",
	Run: func(cmd *cobra.Command, args []string) {
		var ids []string
		for _, p := range strings.Split(socialDMUsers, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				ids = append(ids, p)
			}
		}
		runTool("social.dm.create", map[string]any{"recipients": ids}, flagYes)
	},
}

// messages search 增强 flags（在 extra 里补，避免改太多）
var msgSearchGuild, msgSearchChannel, msgSearchAuthor string
var msgSearchLimit int

func initSocialCommands() {
	stickersPacksCreateCmd.Flags().StringVar(&stickerPackName, "name", "", "名称")
	_ = stickersPacksCreateCmd.MarkFlagRequired("name")
	stickersPacksCmd.AddCommand(stickersPacksListCmd, stickersPacksCreateCmd)
	stickersLibraryCmd.AddCommand(stickersLibraryListCmd)
	stickersAvailableCmd.Flags().StringVar(&stickerGuild, "guild", "", "服务器过滤")
	stickersCmd.AddCommand(stickersPacksCmd, stickersLibraryCmd, stickersAvailableCmd)

	socialPrivacyCmd.AddCommand(socialPrivacyGetCmd)
	socialFriendsRequestCmd.Flags().StringVar(&socialUsername, "username", "", "用户名")
	socialFriendsRequestCmd.Flags().StringVar(&socialUserID, "user", "", "用户 ID")
	socialFriendsAcceptCmd.Flags().StringVar(&socialUserID, "user", "", "用户 ID")
	_ = socialFriendsAcceptCmd.MarkFlagRequired("user")
	socialFriendsBlockCmd.Flags().StringVar(&socialUserID, "user", "", "用户 ID")
	_ = socialFriendsBlockCmd.MarkFlagRequired("user")
	socialFriendsRemoveCmd.Flags().StringVar(&socialUserID, "user", "", "用户 ID")
	_ = socialFriendsRemoveCmd.MarkFlagRequired("user")
	socialFriendsCmd.AddCommand(socialFriendsListCmd, socialFriendsRequestCmd, socialFriendsAcceptCmd, socialFriendsBlockCmd, socialFriendsRemoveCmd)
	socialNotifCmd.AddCommand(socialNotifListCmd)
	socialDMCreateCmd.Flags().StringVar(&socialDMUsers, "users", "", "逗号分隔用户 ID")
	_ = socialDMCreateCmd.MarkFlagRequired("users")
	socialDMCmd.AddCommand(socialDMListCmd, socialDMCreateCmd)
	socialCmd.AddCommand(socialPrivacyCmd, socialFriendsCmd, socialNotifCmd, socialDMCmd)

	// 增强 messages search flags
	messagesSearchCmd.Flags().StringVar(&msgSearchGuild, "guild", "", "限定服务器")
	messagesSearchCmd.Flags().StringVar(&msgSearchChannel, "channel", "", "限定频道")
	messagesSearchCmd.Flags().StringVar(&msgSearchAuthor, "author", "", "作者用户 ID")
	messagesSearchCmd.Flags().IntVar(&msgSearchLimit, "limit", 0, "条数")
	// 覆盖 Run 以带上新参数
	messagesSearchCmd.Run = func(cmd *cobra.Command, args []string) {
		m := map[string]any{"q": msgSearchQ}
		if msgSearchGuild != "" {
			m["guild_id"] = msgSearchGuild
		}
		if msgSearchChannel != "" {
			m["channel_id"] = msgSearchChannel
		}
		if msgSearchAuthor != "" {
			m["author_id"] = msgSearchAuthor
		}
		if msgSearchLimit > 0 {
			m["limit"] = msgSearchLimit
		}
		runTool("messages.search", m, flagYes)
	}
}
