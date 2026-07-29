# 贡献指南

## 开发

```bash
go test ./...
go build -o bin/owl ./cmd/owl
./bin/owl doctor
```

## 结构

| 路径 | 说明 |
|------|------|
| `cmd/owl` | 入口 |
| `internal/tools` | CLI/MCP 共用工具注册表 |
| `internal/mcp` | MCP stdio server |
| `internal/api` | HTTP / OAuth 客户端 |
| `skills/` | AI skill 文档 |

新增工具：在 `internal/tools` 注册，并尽量加 CLI 子命令；MCP 自动暴露。

## 发布

打 tag 触发 Release workflow：

```bash
git tag v0.4.1
git push origin v0.4.1
```

## 与 NewtSpeak 工作区

本目录可独立 git 仓库，也可作为 monorepo 子树。独立发布时仅需 Go 1.22+ 与 Newt-Server OAuth 端点。
