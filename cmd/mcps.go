/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"lain-cli/tools"
	"strings"

	"github.com/spf13/cobra"
)

// mcpsCmd represents the mcps command
var mcpsCmd = &cobra.Command{
	Use:   "mcps",
	Short: "A easy call mcp tools",
	Long: `mcps — 管理并调用已配置的 MCP 服务工具。

用法示例：
  # 列出当前已连接的 MCP 服务
  ./lain-cli mcps

  # 列出某个 MCP 的可用工具
  ./lain-cli mcps <mcp-name>
  例如：
    ./lain-cli mcps rss-reader-mcp

  # 调用某个 MCP 的工具，后续参数为 key:value 形式（会被当作字符串传入）
  ./lain-cli mcps <mcp-name> <tool-name> key1:val1 key2:val2
  例如：
    ./lain-cli mcps whois whois_domain domain:example.com
    ./lain-cli mcps mcp-chinese-fortune fortune-teller year:2005 month:8 day:4 hour:13
`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("=====================", args)
		tools.Init()
		ctx := context.Background()
		tools.ListMCPs()
		// for _, name := range _mcps {
		// 	tools.ListMCPTools(ctx, name)
		// }
		if len(args) == 1 {
			tools.ListMCPTools(ctx, args[0])
		}
		// key:val
		if len(args) >= 2 {

			val := make(map[string]any)

			for i := 2; i < len(args); i++ {

				vals := strings.Split(args[i], ":")
				val[vals[0]] = vals[1]
			}

			tools.CallTool(
				ctx,
				args[0],
				args[1],
				val,
			)
		}
	},
}

func init() {
	rootCmd.AddCommand(mcpsCmd)
}
