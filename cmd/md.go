/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	mui "lain-cli/ui"
	"os"

	"github.com/spf13/cobra"
)

var mdCmd = &cobra.Command{
	Use:   "md",
	Short: "rendering markdown		# 在 terminal 查看markdown 更好的方式",
	Long: `
lain md filename.md
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("file not find")
			return
		}
		filename := args[0]
		if filename == "" {
			fmt.Println("file not find")
			return
		}
		file, err := os.ReadFile(filename)
		if err != nil {
			fmt.Println("open file err")
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
