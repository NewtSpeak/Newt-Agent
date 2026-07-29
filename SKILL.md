---
name: owlspeak-agent
description: Operate NewtSpeak via CLI and MCP — multi-profile OAuth, guild admin, stickers, social, voice, audit, platform; MCP tools/resources/prompts.
---

# NewtSpeak Agent Skill

## 安装与多 profile 登录

```bash
cd Newt-Agent && go install ./cmd/owl

# 设备码（默认）：打开 Desktop/Web 或 owlspeak:// 深链
owl login --server https://a.example --profile home

# PKCE：本机 loopback 回调（适合有浏览器环境）
owl login --server https://a.example --method pkce --client-origin https://app.example

owl login --server https://b.example --profile work --platform
owl profile list && owl profile use work && owl whoami
```

Token 优先 OS keyring；`OWL_AGENT_NO_KEYRING=1` 强制文件。  
用户可在 Desktop **设置 → 已授权应用** 中吊销会话。

## MCP

```bash
owl mcp serve
```

示例配置见 `examples/mcp.json`。

能力：

- **tools/** — 80+ 管理工具  
- **resources/** — `owlspeak://status|tools|whoami|guilds`  
- **prompts/** — `moderate-guild` / `audit-review` / `safe-ops`  

危险工具必须 `confirm: true`。

排障：`owl doctor`（检查 healthz / OAuth / whoami / gapi）。

## 命令域

| 域 | 示例 |
|----|------|
| profile | `owl profile use/list/show/delete` |
| guilds/channels/roles/members | 服管 |
| messages | 含高级 search |
| restrictions/audit/voice | 治理 |
| stickers | 贴图包/库 |
| social | 好友/隐私/通知/私信 |
| platform | 需 `--platform` 登录 |
| tools / mcp | 统一入口 |

`owl tools list` 查看全部。

## 安全

禁止索要密码/refresh；destructive 须确认；权限受 RBAC 约束。
