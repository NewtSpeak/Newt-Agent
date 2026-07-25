package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// 构建时可用 -ldflags "-X github.com/OwlSpeak/Owl-Agent/internal/cmd.Version=..." 覆盖。
var (
	Version   = "0.4.1"
	Commit    = "dev"
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示 owl 版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("owl %s (%s) built %s\n", Version, Commit, BuildDate)
		fmt.Printf("go: %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
