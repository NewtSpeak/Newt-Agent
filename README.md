# Newt-Agent

NewtSpeak **官方 Agent CLI**：以 **用户 OAuth 委托** 管理社区，并作为 **MCP / Skill** 入口供 AI 宿主调用。

```text
newt login (设备码 / PKCE)
        │
        ▼
  access token（keyring）
        │
   ┌────┴────┐
   ▼         ▼
 owl 命令   newt mcp serve  ──► AI（Cursor / Claude Desktop …）
   │         │
   └────┬────┘
        ▼
  tools.Registry  ──Bearer──►  /gapi/v1（业务）
                              /api/v1（平台，需 system_admin）
```

- **不收集密码**：授权在 Desktop / Web 完成  
- 与 [NewtBotSdk](https://github.com/NewtSpeak/NewtBotSdk) 不同：此处是 **真人身份**，不是 bot token  

## 功能

| 能力 | 说明 |
|------|------|
| **多 Profile** | 多 Server / 多账号隔离（`newt profile`） |
| **服管** | 服务器、频道、角色、成员、邀请、消息 |
| **治理** | Restriction、审计、语音踢人/静音/节点池 |
| **社交** | 好友、隐私、通知、私信 |
| **贴图** | 贴图包与库 |
| **平台** | 用户目录、注册开关、SFU 节点（`--platform` 登录） |
| **统一 Tools** | CLI `tools call` 与 MCP **同一注册表**（80+） |
| **MCP** | tools / resources / prompts（stdio JSON-RPC） |
| **Skill** | 给 AI Agent 的操作说明（根目录与 `skills/`） |
| **深链** | `newtspeak://oauth/device` 等配合 Desktop |

## 仓库结构

```text
Newt-Agent/
├── cmd/newt/              # 入口
├── internal/
│   ├── api/              # OAuth + gapi/api HTTP 客户端
│   ├── auth/             # keyring / 会话
│   ├── cmd/              # cobra 子命令
│   ├── config/           # ~/.config/newt-agent
│   ├── mcp/              # MCP server
│   └── tools/            # 工具注册与实现（CLI/MCP 共用）
├── skills/               # 专项 Skill
├── examples/             # mcp.json 等
├── docs/DEEP-LINK.md
├── SKILL.md              # 总 Skill
└── Makefile
```

## 安装与构建

```bash
cd Newt-Agent
make build          # → bin/newt
make test
make install        # go install ./cmd/newt

# 或
go build -o newt ./cmd/newt
./newt doctor
```

发布：打 `v*` tag，CI 交叉编译并上传 GitHub Release。

## 快速开始

```bash
# 设备码（默认）+ 深链
newt login --server https://newt-panel.example.com

# PKCE
newt login --server https://newt-panel.example.com \
  --method pkce --client-origin https://newt-panel.example.com

# 平台管理 scope
newt login --server https://… --platform

# 多 profile
newt login --server https://a.example --profile home
newt profile use home
newt whoami

# 日常
newt guilds list
newt channels list --guild <id>
owl messages send --channel <id> --content "hi"
newt tools list
newt tools call guilds.list --args '{}'
```

危险操作：`--yes` 或 tool 参数 `confirm: true`。

### MCP

```bash
newt mcp serve
```

宿主配置示例：[`examples/mcp.json`](examples/mcp.json)

| 类型 | 内容 |
|------|------|
| tools | 与 `newt tools list` 同源 |
| resources | `newtspeak://status` · `tools` · `whoami` · `guilds` |
| prompts | `moderate-guild` · `audit-review` · `safe-ops` |

## 主要命令

| 域 | 命令 |
|----|------|
| 会话 | `login` `logout` `whoami` `status` `doctor` `profile` |
| 服管 | `guilds` `channels` `roles` `members` `invites` `messages` |
| 治理 | `restrictions` `audit` `voice` |
| 扩展 | `stickers` `social` `platform` |
| 集成 | `tools` `mcp serve` |

## 文档

| 文档 | 内容 |
|------|------|
| [**docs/USAGE.md**](docs/USAGE.md) | **使用指南**（推荐入口） |
| [**docs/API.md**](docs/API.md) | 命令树、Tools 全表、OAuth/gapi |
| [docs/DEEP-LINK.md](docs/DEEP-LINK.md) | 深链与单实例 |
| [SKILL.md](SKILL.md) | AI Skill 总说明 |
| [skills/newtspeak-guild-admin/](skills/newtspeak-guild-admin/) | 服管专项 Skill |
| [CONTRIBUTING.md](CONTRIBUTING.md) | 贡献指南 |
| [examples/mcp.json](examples/mcp.json) | MCP 配置示例 |

## 安全

1. 勿向 AI 提供密码或 refresh token  
2. Token 默认 OS keyring（`NEWT_AGENT_NO_KEYRING=1` 可改文件）  
3. 危险 tool 必须确认  
4. 权限完全受服内 RBAC + OAuth scope 约束  
5. 可在 Desktop **设置 → 已授权应用** 吊销  

## 相关仓库

| 仓库 | 关系 |
|------|------|
| [Newt-Server](https://github.com/NewtSpeak/Newt-Server) | OAuth、`/gapi/v1`、`/api/v1` |
| [Newt-Desktop](https://github.com/NewtSpeak/Newt-Desktop) | 授权页与深链 |
| [NewtBotSdk](https://github.com/NewtSpeak/NewtBotSdk) | 机器人 SDK（另一身份平面） |

## 许可证

见仓库 [`LICENSE`](./LICENSE)。
