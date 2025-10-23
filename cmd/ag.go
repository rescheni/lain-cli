/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"lain-cli/server"
	"lain-cli/utils"
	"os"

	"github.com/spf13/cobra"
)

var (
	flagMarkdown bool = false
)

// agCmd represents the ag command
var agCmd = &cobra.Command{
	Use:   "ag",
	Short: "input Your things to lain",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		ask := ""
		if len(args) > 0 {
			ask += "line input:"
			ask += fmt.Sprintln(args)
		}
		// fmt.Println(ask)

		stat, _ := os.Stdin.Stat()
		ask += "pipe input:"
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				ask += scanner.Text()
			}
		}

		// fmt.Println(ask)
		ask = utils.Prompt + ask
		if flagMarkdown {
			server.CallModel(ctx, ask)
		} else {
			err := server.CallModelStream(ctx, ask)
			if err != nil {
				// fmt.Println("place input things to llm")
				fmt.Printf("server call model error")
			}
		}

		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(agCmd)
	agCmd.Flags().BoolVarP(&flagMarkdown, "markdown", "m", false, "Print Use Markdown")
}
