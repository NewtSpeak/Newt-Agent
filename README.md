# Owl-Agent

OwlSpeak **CLI + MCP + Skill**（OAuth 用户委托）。

可独立成仓库；CI 在 push/PR 上测试并交叉编译，打 `v*` tag 自动发 GitHub Release。

## 构建 / 测试

```bash
make build          # → bin/owl
make test
./bin/owl doctor

# 或
go build -o owl ./cmd/owl
go test ./...
```

发布（维护者）：

```bash
git tag v0.4.1 && git push origin v0.4.1
```

MCP 配置：`examples/mcp.json`  
深链：`docs/DEEP-LINK.md`  
贡献：`CONTRIBUTING.md`

## 登录方式

```bash
# 设备码（默认）+ 深链 owlspeak://oauth/device
owl login --server https://api.example

# PKCE loopback
owl login --server https://api.example --method pkce --client-origin https://web.example
```

## 多 Profile

```bash
owl login --server https://a.example --profile home
owl login --server https://b.example --profile work
owl profile list
owl profile use work
owl profile show
```

## MCP

- tools / resources / prompts（v0.4.0）
- resources: `owlspeak://status` `owlspeak://tools` `owlspeak://whoami` `owlspeak://guilds`
- prompts: `moderate-guild` `audit-review` `safe-ops`

## 主要命令

`guilds` `channels` `roles` `members` `invites` `messages`  
`restrictions` `audit` `voice` `stickers` `social` `platform`  
`profile` `tools` `mcp serve`
