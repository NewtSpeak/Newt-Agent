---
name: owlspeak-guild-admin
description: Full NewtSpeak guild administration via owl CLI/MCP — channels, roles, members, invites, messages.
---

# NewtSpeak 服务器管理

依赖 `owl login` 与根 skill `owlspeak-agent`。

## 频道

```bash
owl channels list --guild <gid>
owl channels create --guild <gid> --name welcome --type TEXT
owl channels update --channel <cid> --topic "欢迎"
owl channels delete --channel <cid> --yes
owl channels overwrites list --guild <gid> --channel <cid>
owl channels overwrites set --channel <cid> --target <role_id> --type ROLE --allow 0 --deny 0
```

## 角色与成员

```bash
owl roles list --guild <gid>
owl roles create --guild <gid> --name Moderator --permissions 0 --position 1
owl roles assign --guild <gid> --member <mid> --role <rid>
owl members list --guild <gid>
owl members kick --guild <gid> --member <mid> --yes
```

## MCP

同一套 tools：`channels.*` `roles.*` `members.*` …  
写操作 arguments 加 `"confirm": true`。
