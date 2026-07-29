package cmd

import (
	"fmt"
	"os"

	"github.com/NewtSpeak/Newt-Agent/internal/mcp"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "MCP（Model Context Protocol）相关",
}

var mcpServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "以 stdio 运行 MCP server（供 AI 客户端连接）",
	Long: `启动 MCP server，通过标准输入/输出与宿主通信。

配置示例（Claude Desktop / Cursor mcp.json）：

  {
    "mcpServers": {
      "newtspeak": {
        "command": "owl",
        "args": ["mcp", "serve"]
      }
    }
  }

使用前请先完成: newt login --server <url>
危险写操作需在 arguments 中传 confirm: true。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 预检登录，失败时明确写到 stderr（stdout 留给 JSON-RPC）
		if _, err := makeClient(); err != nil {
			fmt.Fprintln(os.Stderr, "newt mcp serve: 未登录或会话无效:", err)
			fmt.Fprintln(os.Stderr, "请先运行: newt login --server <url>")
			// 仍启动 server：部分客户端先 initialize 再 tools；调用时再报错
		}
		s := &mcp.Server{Reg: toolRegistry()}
		if err := s.Serve(); err != nil {
			fmt.Fprintln(os.Stderr, "mcp serve 退出:", err)
			os.Exit(1)
		}
	},
}

func init() {
	mcpCmd.AddCommand(mcpServeCmd)
}
