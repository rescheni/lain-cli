/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"lain-cli/tools"

	"github.com/spf13/cobra"
)

// topCmd represents the top command
var topCmd = &cobra.Command{
	Use:   "top",
	Short: "Open Tui View ",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		tools.OpenPerformance()
	},
}

func init() {
	rootCmd.AddCommand(topCmd)
}
