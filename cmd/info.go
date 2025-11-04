/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"lain-cli/tools"

	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "show device info",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		tools.InfoInit()
		tools.BasePrint()
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
