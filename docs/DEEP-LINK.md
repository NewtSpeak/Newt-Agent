# 深链与单实例

## Scheme

| URL | 作用 |
|-----|------|
| `owlspeak://oauth/device?user_code=XXXX-XXXX&server=https://...` | 打开设备码授权页 |
| `owlspeak://oauth/authorize?...` | PKCE 授权页 |
| `owlspeak://invite?code=...&server=...` | 社区邀请 |
| `owlspeak://register?code=...&server=...` | 注册邀请 |

## 桌面端实现

1. **tauri-plugin-deep-link**：注册 `owlspeak` scheme（`tauri.conf.json` → plugins.deep-link）。
2. **tauri-plugin-single-instance**：应用已运行时，二次启动/深链把 argv 中的 URL 经事件 `owl://deep-link` 发给已有窗口并 `set_focus`。
3. **前端** `useDeepLinkNavigation`：同时订阅 `onOpenUrl` 与 `owl://deep-link`。

## CLI

```bash
owl login --server https://api.example
# 会尝试打开浏览器 + owlspeak://oauth/device?...
```

开发态：`tauri dev` 一般会注册 scheme；正式安装包安装后 OS 级注册生效。
