# Owl-Agent 接口与调用文档

三层入口同一工具注册表：

```text
owl <命令>  ──┐
owl tools call ─┼──► tools.Registry ──► /gapi/v1 或 /api/v1
owl mcp serve  ─┘         │
                          └── Bearer <access_token>
```

## 认证平面（OAuth）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/oauth/v1/device/code` | 设备码 |
| POST | `/oauth/v1/token` | device_code / refresh / PKCE |
| POST | `/oauth/v1/revoke` | 吊销 |
| GET | `/oauth/v1/userinfo` | 当前用户 |

登录后业务请求：

```http
Authorization: Bearer <access_token>
```

| 前缀 | 用途 | 条件 |
|------|------|------|
| `/gapi/v1` | 用户业务 API（服/频道/消息…） | scope 含 `gapi.full` |
| `/api/v1` | 平台管理 | `platform.*` + `system_admin` |

Agent 内部：`Client.Gapi(method, path, body, query)` / `Client.Api(...)`。

---

## CLI 命令树

全局：`--yes` 跳过危险确认 · `owl completion <shell>`

### 会话

| 命令 | 说明 |
|------|------|
| `owl login --server URL [--method device\|pkce] [--platform] [--profile P] [--no-open]` | 登录 |
| `owl logout` | 登出 |
| `owl whoami` | 当前用户 |
| `owl status` | 本地配置与 token 后端 |
| `owl doctor` | 连通性诊断 |
| `owl profile list\|use\|show\|delete` | 多 profile |

### 服管

| 命令 | 对应 tool |
|------|-----------|
| `owl guilds list\|get\|create\|update\|delete\|permissions` | `guilds.*` |
| `owl channels list\|create\|update\|delete\|…` | `channels.*` |
| `owl roles list\|create\|…` | `roles.*` |
| `owl members list\|kick\|ban\|…` | `members.*` |
| `owl invites create\|…` | `invites.*` |
| `owl messages list\|get\|send\|edit\|delete\|search` | `messages.*` |

### 治理 / 语音

| 命令 | tool |
|------|------|
| `owl restrictions list\|create\|lift` | `restrictions.*` |
| `owl audit list` | `audit.list` |
| `owl voice states\|disconnect\|move\|mute\|nodes\|…` | `voice.*` |

### 贴图 / 社交 / 平台

| 命令 | tool |
|------|------|
| `owl stickers packs\|library\|available` | `stickers.*` |
| `owl social privacy\|friends\|notifications\|dm` | `social.*` |
| `owl platform …` | `platform.*` |

### 工具与 MCP

| 命令 | 说明 |
|------|------|
| `owl tools list` | 列出全部 tools（JSON） |
| `owl tools call <name> --args '{…}' [--yes]` | 直接调 tool |
| `owl mcp serve` | MCP stdio 服务 |

示例：

```bash
owl channels list --guild 018f…
owl messages send --channel 018f… --content "hello"
owl tools call roles.assign --args '{
  "guild_id":"…","user_id":"…","role_id":"…"
}'
owl members ban --guild … --user … --yes
```

---

## Tools 目录（CLI / MCP 共用）

危险工具标 **⚠**，调用须 `confirm: true` 或 CLI `--yes`。

### 元信息

| Name | 说明 |
|------|------|
| `whoami` | OAuth 用户信息 |
| `status` | 本地登录状态（可不登录） |

### 服务器 `guilds.*`

| Name | 主要参数 | ⚠ |
|------|----------|---|
| `guilds.list` | — | |
| `guilds.get` | `guild_id` | |
| `guilds.create` | `name` | |
| `guilds.update` | `guild_id`, name/description/… | |
| `guilds.delete` | `guild_id`, `confirm_name` | ⚠ |
| `guilds.permissions` | `guild_id` | |
| `guilds.transfer` | `guild_id`, 新所有者 | ⚠ |

### 频道 `channels.*`

| Name | 说明 | ⚠ |
|------|------|---|
| `channels.list` | `guild_id` | |
| `channels.create` | `guild_id`, `name`, `type`=`TEXT\|VOICE\|CATEGORY\|STAGE` | |
| `channels.update` | `channel_id`, … | |
| `channels.delete` | `channel_id` | ⚠ |
| `channels.reorder` | `guild_id`, `items[]` | |
| `channels.permissions` | `channel_id` | |
| `channels.overwrites.list` | `channel_id` | |
| `channels.overwrites.upsert` | target + allow/deny | |
| `channels.overwrites.delete` | | ⚠ |

### 角色 `roles.*`

| Name | 说明 | ⚠ |
|------|------|---|
| `roles.list` / `create` / `update` | | |
| `roles.delete` | | ⚠ |
| `roles.assign` / `remove` | `guild_id`, `user_id`, `role_id` | |

### 成员 `members.*`

| Name | 说明 | ⚠ |
|------|------|---|
| `members.list` | | |
| `members.nick` | | |
| `members.kick` | `@me` = 退服 | ⚠ |
| `members.ban` / `unban` | | ⚠ |
| `members.bans` | 封禁列表 | |

### 邀请 / 消息

| Name | 说明 | ⚠ |
|------|------|---|
| `invites.create` / `get` / `join` | | |
| `messages.list` | `channel_id`, before/after/limit | |
| `messages.get` | | |
| `messages.send` | `channel_id`, `content` | |
| `messages.edit` | | |
| `messages.delete` | | ⚠ |
| `messages.search` | `q` + guild/channel/author/… | |

### 限制 / 审计 / 语音

| Name | 说明 | ⚠ |
|------|------|---|
| `restrictions.list` / `get` / `patch` | | |
| `restrictions.create` | scope/kind/deny | ⚠ |
| `restrictions.lift` | | ⚠ |
| `audit.list` | `guild_id`, action, limit | |
| `voice.states` | | |
| `voice.disconnect` / `move` / `server_mute` | | ⚠ |
| `voice.nodes` | 候选节点 | |
| `voice.node_pool.get` / `set` | 服节点池 | |

Restriction `scope`：`TEXT_CHANNEL` / `VOICE_CHANNEL` / `GUILD_ALL_TEXT` / `GUILD_ALL_VOICE`；`kind`：`SANCTION` / `CHANNEL_BAN`。

### 贴图 `stickers.*`

| Name | 说明 |
|------|------|
| `stickers.packs.list` / `get` / `create` / `delete` | 自建包 |
| `stickers.library.list` / `install` / `uninstall` | 库 |
| `stickers.available` | 可用集合（可选 `guild_id`） |
| `stickers.guild_bans.list` / `add` / `remove` | 服 ban |

### 社交 `social.*`

| Name | 说明 | ⚠ |
|------|------|---|
| `social.privacy.get` / `patch` | 隐私 | |
| `social.friends.list` / `request` / `accept` | 好友 | |
| `social.friends.block` / `remove` | | ⚠ |
| `social.notifications.list` / `ack` | 通知 | |
| `social.dm.list` / `create` | 私信/群 | |

### 平台 `platform.*`（需 platform scope）

| Name | 说明 | ⚠ |
|------|------|---|
| `platform.users.list` | 用户目录 | |
| `platform.users.disable` / `enable` | | ⚠ |
| `platform.users.reset_password` | | ⚠ |
| `platform.users.set_admin` | | ⚠ |
| `platform.registration.get` / `set` | 注册开关 | |
| `platform.sfu.nodes` / `topology` | SFU | |
| `platform.audit.list` | 全站审计 | |

完整 schema：`owl tools list` 或读 `Owl-Agent/internal/tools/handlers*.go` 的 `InputSchema`。

---

## tools.call 约定

```bash
owl tools call <name> --args '<JSON object>'
```

```json
{
  "guild_id": "018f…",
  "confirm": true
}
```

| 规则 | 说明 |
|------|------|
| 大 ID 用**字符串** | JSON number 会丢雪花精度 |
| 危险工具 | 必须 `confirm: true`，除非 `--yes` / `force` |
| `status` | 唯一可不登录的 tool |
| 返回 | JSON（成功对象或错误信息） |

MCP 侧等价：`tools/call`，arguments 同 schema；destructive 同样要 `confirm`。

---

## MCP 协议摘要

```text
传输: stdio，每行一条 JSON-RPC 2.0
协议: 2024-11-05
server: owl-agent
```

| 方法 | 作用 |
|------|------|
| `initialize` | 握手 |
| `tools/list` · `tools/call` | 工具 |
| `resources/list` · `resources/read` | 资源 |
| `prompts/list` · `prompts/get` | 提示词 |

**Resources**

| URI | 内容 |
|-----|------|
| `owlspeak://status` | 登录/配置 |
| `owlspeak://tools` | 工具清单 |
| `owlspeak://whoami` | 用户信息 |
| `owlspeak://guilds` | 已加入服务器 |

**Prompts**：`moderate-guild` · `audit-review` · `safe-ops`

配置样例：`Owl-Agent/examples/mcp.json`、`examples/claude_desktop_mcp.json`。

---

## 底层 gapi 路径映射（节选）

Tools 最终打到用户 API（路径随 Server 版本微调，以源码为准）：

| Tool 域 | 典型路径前缀 |
|---------|----------------|
| guilds | `/gapi/v1/guilds` |
| channels | `/gapi/v1/guilds/{id}/channels`、`/gapi/v1/channels/{id}` |
| roles | `/gapi/v1/guilds/{id}/roles` |
| members | `/gapi/v1/guilds/{id}/members` |
| messages | `/gapi/v1/channels/{id}/messages` |
| invites | `/gapi/v1/guilds/{id}/invites`、`/gapi/v1/invites/{code}` |
| restrictions | `/gapi/v1/guilds/{id}/restrictions` |
| audit | `/gapi/v1/guilds/{id}/audit-logs` |
| voice | `/gapi/v1/guilds/{id}/voice/…` |
| stickers | `/gapi/v1/users/@me/sticker-…` |
| social | `/gapi/v1/users/@me/…`（好友/隐私/DM） |
| platform | `/api/v1/…`（平台） |

直接 HTTP 调试（已登录可从 keyring 取 token，**勿提交**）：

```bash
curl -sH "Authorization: Bearer $TOKEN" \
  https://owl-panel.example.com/gapi/v1/guilds
```

---

## 错误

| 来源 | 表现 |
|------|------|
| HTTP 4xx/5xx | `CODE: message`（解析 `error.code`） |
| 未知 tool | 列出可用 name |
| 缺 confirm | `危险操作 xxx 需要 confirm=true` |
| 未登录 | factory 取 session 失败 → 先 `owl login` |

---

## 安全清单

1. 不向 AI 提供密码、refresh token、keyring 导出  
2. destructive 必须二次确认  
3. 权限完全受服内 RBAC + OAuth scope 约束  
4. 平台工具仅 system_admin  
5. 吊销：Desktop 已授权应用，或 `owl logout` + `/oauth/v1/revoke`
