/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"context"
	"fmt"

	"github.com/rescheni/lain-cli/config"
	"github.com/rescheni/lain-cli/internal/server"
	"github.com/rescheni/lain-cli/internal/tools"
	"github.com/rescheni/lain-cli/logs"

	"os"

	"github.com/spf13/cobra"
)

var (
	flagMarkdown bool = false
	flagWindow   bool = false
)

// agCmd represents the ag command
var agCmd = &cobra.Command{
	Use:   "ag",
	Short: "input Your things to lain",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		logs.Debug(os.TempDir())
		err := server.LLMInit()
		if err != nil {
			logs.Err("server LLM init Error")
		}

		if len(args) == 0 || len(args) > 1 && args[0] == "" {
			logs.Info("Please enter a request")
			return
		}
		ctx := context.Background()
		ask := ""
		if len(args) > 0 {
			ask += "line input:"
			ask += fmt.Sprintln(args)
		}
		if config.Conf.Context.Enabled {
			tools.LLMCTX.Init()
		}

		stat, _ := os.Stdin.Stat()
		ask += "pipe input:"
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				ask += scanner.Text()
			}
		}
		if config.Conf.Context.Enabled {
			tools.LLMCTX.Add(ask + "\n")
			ask = tools.LLMCTX.Getcontext() + ask

		}

		// Use markdown out
		if flagMarkdown {
			server.CallModel(ctx, ask, flagWindow)
		} else {
			err := server.CallModelStream(ctx, ask)
			if err != nil {
				logs.Err("server call model error", err)
			}
		}
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(agCmd)
	agCmd.Flags().BoolVarP(&flagMarkdown, "markdown", "m", false, "Print Use Markdown	# 使用 markdown 格式输出")
	agCmd.Flags().BoolVarP(&flagWindow, "window", "w", false, "print markdown on new window		# 新窗口显示")
}
