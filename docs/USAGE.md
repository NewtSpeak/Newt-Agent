# Newt-Agent（CLI）使用文档

用户委托的 **CLI + MCP + Skill**：OAuth 登录后，以真人身份操作 NewtSpeak（服管、治理、社交、平台）。

> 命令与工具全表见 [API.md](./API.md) · 深链 [DEEP-LINK.md](./DEEP-LINK.md)

## 与 Bot SDK 的区别

| | Agent CLI | Bot SDK |
|--|-----------|---------|
| 身份 | OAuth 用户 | 机器人 |
| 凭证 | access/refresh（keyring） | `newtbot_…` |
| API | `/gapi/v1`、可选 `/api/v1` | `/bot-api/v1` |
| 入口 | `owl` 命令 / MCP / Skill | 各语言库 |

**密码不进 CLI**：设备码或 PKCE 在 Desktop/Web 授权页完成。

## 1. 安装

```bash
# 源码
cd Newt-Agent
make build          # → bin/owl
make install        # go install ./cmd/owl

# 或
go build -o owl ./cmd/owl
```

Release：各平台 `owl` 二进制（版本见 Makefile，当前约 `0.4.x`）。

```bash
owl doctor          # healthz / OAuth / whoami / gapi
owl version
```

Shell 补全：

```bash
owl completion bash > /etc/bash_completion.d/owl
owl completion powershell | Out-String | Invoke-Expression
```

## 2. 登录

### 设备码（默认）

```bash
owl login --server https://newt-panel.example.com
# 打印 user_code + 打开浏览器 / newtspeak://oauth/device
```

### PKCE（本机 loopback）

```bash
owl login --server https://newt-panel.example.com \
  --method pkce \
  --client-origin https://newt-panel.example.com
```

### 平台管理 scope

```bash
owl login --server https://… --platform
# scope 追加 platform.read platform.admin（需 system_admin）
```

### 多 Profile

```bash
owl login --server https://a.example --profile home
owl login --server https://b.example --profile work
owl profile list
owl profile use work
owl profile show
owl whoami
owl logout
```

| 项 | 说明 |
|----|------|
| 配置目录 | `~/.config/owl-agent/`（Windows: `%AppData%\owl-agent\`） |
| Token | 优先 OS keyring；`OWL_AGENT_NO_KEYRING=1` 强制文件 |
| 吊销 | Desktop **设置 → 已授权应用** |

## 3. 日常命令速览

```bash
# 服务器
owl guilds list
owl guilds get <gid>
owl guilds create --name "测试服"
owl guilds permissions <gid>

# 频道
owl channels list --guild <gid>
owl channels create --guild <gid> --name general --type TEXT

# 角色 / 成员
owl roles list --guild <gid>
owl members list --guild <gid>
owl members kick --guild <gid> --user <uid> --yes

# 消息
owl messages list --channel <cid>
owl messages send --channel <cid> --content "hi"
owl messages search --guild <gid> --q "关键词"

# 邀请 / 限制 / 审计 / 语音
owl invites create --guild <gid>
owl restrictions list --guild <gid>
owl audit list --guild <gid>
owl voice states --guild <gid> --channel <cid>

# 贴图 / 社交
owl stickers packs list
owl social friends list
owl social dm list

# 平台（需 --platform 登录）
owl platform users list
owl platform sfu nodes
```

危险操作加 `--yes`（等同 tool 的 `confirm=true`）。

完整子命令见 [API.md](./API.md)。

## 4. 统一工具入口

CLI 与 MCP **共用同一套 tools**：

```bash
owl tools list                          # JSON：name / description / destructive
owl tools call guilds.list --args '{}'
owl tools call messages.send --args '{"channel_id":"…","content":"hi"}'
owl tools call members.kick --args '{"guild_id":"…","user_id":"…","confirm":true}'
# 或
owl tools call members.kick --yes --args '{"guild_id":"…","user_id":"…"}'
```

## 5. MCP（给 AI 宿主）

```bash
owl mcp serve          # stdio JSON-RPC 2.0
```

**Cursor / Claude Desktop 示例**（`examples/mcp.json`）：

```json
{
  "mcpServers": {
    "newtspeak": {
      "command": "owl",
      "args": ["mcp", "serve"],
      "env": {}
    }
  }
}
```

| 能力 | 内容 |
|------|------|
| **tools** | 80+（与 `owl tools list` 同源） |
| **resources** | `newtspeak://status` `tools` `whoami` `guilds` |
| **prompts** | `moderate-guild` / `audit-review` / `safe-ops` |

危险工具参数必须带 `confirm: true`（或宿主侧等价确认）。

## 6. Skill

- 根目录 `Newt-Agent/SKILL.md`：总 skill  
- `skills/newtspeak-guild-admin/SKILL.md`：服管专项  

AI Agent 加载 skill 后，优先通过 MCP/`owl tools call` 操作，**禁止索要密码或 refresh token**。

## 7. 配置与环境

| 变量 / 项 | 说明 |
|-----------|------|
| `--server` | Newt-Server 公网根（无尾斜杠） |
| `--profile` | 多账号隔离 |
| `OWL_AGENT_NO_KEYRING=1` | token 写入 config 文件 |
| config.json | `server_url` / `scope` / `client_id` 等 |

默认 OAuth scope：`openid profile gapi.full offline_access`。

## 8. 排障

```bash
owl doctor
owl status
owl whoami
```

| 现象 | 处理 |
|------|------|
| 未登录 | `owl login --server …` |
| 403 权限 | 账号在服内 RBAC 不足 |
| platform.* 失败 | 用 `--platform` 重登 + system_admin |
| MCP 无 tools | 确认 `owl` 在 PATH；先 CLI 登录同一 profile |
| 设备码超时 | 重新 `login`；检查 Desktop 深链 |

## 9. 相关文档

| 文档 | 说明 |
|------|------|
| [API.md](./API.md) | 命令、Tools、底层 HTTP |
| [NewtBotSdk](https://github.com/NewtSpeak/NewtBotSdk) | 机器人 SDK |
| [Server 部署](https://github.com/NewtSpeak/Newt-Server/blob/main/docs/deploy/server.md) | Server 部署 |
