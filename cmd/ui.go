/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/spf13/cobra"
)

var (
	purple    = lipgloss.Color("99")
	gray      = lipgloss.Color("245")
	lightGray = lipgloss.Color("241")

	headerStyle  = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
	cellStyle    = lipgloss.NewStyle().Padding(0, 1).Width(14)
	oddRowStyle  = cellStyle.Foreground(gray)
	evenRowStyle = cellStyle.Foreground(lightGray)
)
var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	PaddingTop(2).
	PaddingLeft(4).
	Width(22)

// uiCmd represents the ui command
var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Open Lain-CLI UI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ui called")
		fmt.Println(style.Render("Hello, kitty"))
		rows := [][]string{
			{"Chinese", "您好", "你好"},
			{"Japanese", "こんにちは", "やあ"},
			{"Arabic", "أهلين", "أهلا"},
			{"Russian", "Здравствуйте", "Привет"},
			{"Spanish", "Hola", "¿Qué tal?"},
		}
		t := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(purple)).
			StyleFunc(func(row, col int) lipgloss.Style {
				switch {
				case row == table.HeaderRow:
					return headerStyle
				case row%2 == 0:
					return evenRowStyle
				default:
					return oddRowStyle
				}
			}).
			Headers("Download", "Upload", "Time").
			Rows(rows...)
		t.Row("English", "You look absolutely fabulous.", "How's it going?")
		fmt.Println(t)

	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
	rootCmd.Flags().BoolP("ui", "u", false, "Open Lain-CLI ui")
}
