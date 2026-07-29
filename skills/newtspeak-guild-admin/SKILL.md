---
name: newtspeak-guild-admin
description: Full NewtSpeak guild administration via owl CLI/MCP — channels, roles, members, invites, messages.
---

# NewtSpeak 服务器管理

依赖 `newt login` 与根 skill `newtspeak-agent`。

## 频道

```bash
newt channels list --guild <gid>
newt channels create --guild <gid> --name welcome --type TEXT
newt channels update --channel <cid> --topic "欢迎"
newt channels delete --channel <cid> --yes
newt channels overwrites list --guild <gid> --channel <cid>
newt channels overwrites set --channel <cid> --target <role_id> --type ROLE --allow 0 --deny 0
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
