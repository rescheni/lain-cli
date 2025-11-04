/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/rescheni/lain-cli/internal/tools"
	"github.com/spf13/cobra"
)

// topCmd represents the top command
var topCmd = &cobra.Command{
	Use:   "top",
	Short: "Open Tui View Process Memory Cpu Status  # 实时查看进程内存Cpu等的信息",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		tools.OpenPerformance()
	},
}

func init() {
	rootCmd.AddCommand(topCmd)
}
