/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/rescheni/lain-cli/config"
	"github.com/rescheni/lain-cli/internal/tools"
	logs "github.com/rescheni/lain-cli/logs"

	"github.com/spf13/cobra"
)

var tofile string

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
		err := tools.Init()
		if err != nil {
			return
		}
		ctx := context.Background()
		_mcps := tools.ListMCPs()
		for i, name := range _mcps {
			fmt.Println(i, name)
		}

		if len(args) == 1 {
			tools.ListMCPTools(ctx, args[0])
		}
		// key:val
		if len(args) >= 2 {

			val := make(map[string]any)

			for i := 2; i < len(args); i++ {

				vals := strings.Split(args[i], "===")
				val[vals[0]] = vals[1]
			}

			tools.CallTool(
				ctx,
				args[0],
				args[1],
				val,
				tofile,
			)
		}
	},
}

// 交互方式调用mcp
var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "call mcp interactive",
	Long:  `交互调用mcp`,
	Run: func(cmd *cobra.Command, args []string) {

		// 捕获 Ctrl+C [TODO] 关闭PY依然大量输出报错
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGINT)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT)

		// 初始化mcp
		tools.Init()
		ctx := context.Background()
		mcps := tools.ListMCPs()
		// ls 		 								-- 列出所有mcp
		// list [mcp] 								-- 列出某个mcp的工具
		// exec [mcp] [mcp tools] [vals]			-- 运行某个mcp 工具
		// 交互开始
		var ok bool
		go func() {
			_, ok = <-sigs
		}()
		for {
			// line := ""
			fmt.Print("Lain-> ")
			// fmt.Scanln(&line)
			reader := bufio.NewReader(os.Stdin)
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			arg := strings.Split(line, " ")
			// fmt.Println(line)
			if len(arg) == 1 && arg[0] == "exit" || ok {
				logs.Info("Exited")
				return
			} else if len(arg) == 1 && arg[0] == "ls" {
				fmt.Println("MCP list:")
				for i, v := range mcps {
					fmt.Printf("\t %d-%s\n", i+1, v)
				}
			} else if len(arg) == 2 && arg[0] == "list" {
				tools.ListMCPTools(ctx, arg[1])
			} else if len(arg) >= 3 && arg[0] == "exec" {
				vals := make(map[string]any)
				for i := 3; i < len(arg); i++ {
					temp := strings.Split(arg[i], "===")
					vals[temp[0]] = temp[1]
				}
				tools.CallTool(ctx, arg[1], arg[2], vals, tofile)
			} else if len(arg) == 0 {
				continue
			} else {
				logs.Err("Not find command")
			}

		}
	},
}

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit mcp.json",
	Long:  `编辑mcp.json文件`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var editorCmd *exec.Cmd
		if len(args) == 0 {
			_, err = exec.LookPath("vim")
			if err != nil {
				_, err = exec.LookPath("nano")
				if err == nil {
					editorCmd = exec.Command("nano", config.Conf.Mcp.Json)
				} else {
					err = errors.New("no vim/nano")
				}
			} else {
				editorCmd = exec.Command("vim", config.Conf.Mcp.Json)

			}
		} else {
			_, err = exec.LookPath(args[0])
			if err != nil {
				err = errors.New("no " + args[0])
			} else {
				editorCmd = exec.Command(args[0], config.Conf.Mcp.Json)
			}
		}
		if err != nil {
			logs.Err("edit file err : ", err)
		} else {
			logs.Info("open " + config.Conf.Mcp.Json + " OK")
			editorCmd.Stdin = os.Stdin
			editorCmd.Stdout = os.Stdout
			editorCmd.Stderr = os.Stderr
			if err = editorCmd.Run(); err != nil {
				logs.Err("err", err)
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(mcpsCmd)
	mcpsCmd.AddCommand(replCmd)
	mcpsCmd.AddCommand(editCmd)
	mcpsCmd.Flags().StringVarP(&tofile, "tofile", "f", "", "Mcp print to file	# 将mcp的输出同时到文件")
}
