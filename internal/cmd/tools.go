package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "列举或调用 AI 工具（与 MCP 同源）",
}

var toolsListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出全部 tools",
	Run: func(cmd *cobra.Command, args []string) {
		reg := toolRegistry()
		type row struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Destructive bool   `json:"destructive,omitempty"`
			CLI         string `json:"cli_hint,omitempty"`
		}
		var out []row
		for _, d := range reg.List() {
			out = append(out, row{d.Name, d.Description, d.Destructive, d.CLIHint})
		}
		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(b))
	},
}

var toolsCallArgs string
var toolsCallYes bool

var toolsCallCmd = &cobra.Command{
	Use:   "call <name>",
	Short: "调用命名 tool（--args JSON）",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := strings.TrimSpace(args[0])
		var parsed map[string]any
		if err := json.Unmarshal([]byte(toolsCallArgs), &parsed); err != nil {
			fatal(fmt.Errorf("--args 必须是 JSON 对象: %w", err))
		}
		runTool(name, parsed, toolsCallYes || flagYes)
	},
}

func init() {
	toolsCallCmd.Flags().StringVar(&toolsCallArgs, "args", "{}", "JSON 参数对象")
	toolsCallCmd.Flags().BoolVar(&toolsCallYes, "yes", false, "危险操作跳过 confirm")
	toolsCmd.AddCommand(toolsListCmd, toolsCallCmd)
}
