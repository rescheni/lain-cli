/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// agCmd represents the ag command
var agCmd = &cobra.Command{
	Use:   "ag",
	Short: "input Your things to lain",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Println(args[0])
			return
		} else {
			fmt.Println("no args")
		}
	},
}

func init() {
	rootCmd.AddCommand(agCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// agCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// agCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
