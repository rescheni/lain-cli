/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	mui "github.com/rescheni/lain-cli/internal/ui"

	"github.com/rescheni/lain-cli/logs"

	"github.com/spf13/cobra"
)

var mdCmd = &cobra.Command{
	Use:   "md",
	Short: "rendering markdown		# 在 terminal 查看markdown 更好的方式",
	Long: `
lain-cli md filename.md
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			logs.Err("file not find")
			return
		}
		filename := args[0]
		if filename == "" {
			logs.Err("file not find")
			return
		}
		file, err := os.ReadFile(filename)
		if err != nil {
			logs.Err("open file err")
			return
		}
		// fmt.Println()
		mui.PrintMarkdown(string(file), flagWindow)

	},
}

func init() {
	rootCmd.AddCommand(mdCmd)
	mdCmd.Flags().BoolVarP(&flagWindow, "window", "w", false, "print markdown on new window		# 新窗口显示")

}
