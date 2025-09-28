package cmd

import (
	"lain-cli/utils"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// 渐变色工具函数
func gradient(text string, colors []string) string {
	lines := strings.Split(text, "\n") // 替代 lipgloss.SplitLines
	styled := make([]string, len(lines))
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue // 空行跳过，避免多余颜色
		}
		color := colors[i%len(colors)]
		styled[i] = lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Bold(true).
			Render(line)
	}
	return lipgloss.JoinVertical(lipgloss.Left, styled...)
}

var rootCmd = &cobra.Command{
	Use:   "lain-cli",
	Short: "A TUI TOOLS",
	Long: func() string {
		logo := `
█████                  ███                          █████████  █████       █████
▒▒███                  ▒▒▒                          ███▒▒▒▒▒███▒▒███       ▒▒███ 
 ▒███         ██████   ████  ████████              ███     ▒▒▒  ▒███        ▒███ 
 ▒███        ▒▒▒▒▒███ ▒▒███ ▒▒███▒▒███  ██████████▒███          ▒███        ▒███ 
 ▒███         ███████  ▒███  ▒███ ▒███ ▒▒▒▒▒▒▒▒▒▒ ▒███          ▒███        ▒███ 
 ▒███      █ ███▒▒███  ▒███  ▒███ ▒███            ▒▒███     ███ ▒███      █ ▒███ 
 ███████████▒▒████████ █████ ████ █████            ▒▒█████████  ███████████ █████
▒▒▒▒▒▒▒▒▒▒▒  ▒▒▒▒▒▒▒▒ ▒▒▒▒▒ ▒▒▒▒ ▒▒▒▒▒              ▒▒▒▒▒▒▒▒▒  ▒▒▒▒▒▒▒▒▒▒▒ ▒▒▒▒▒
`

		// 定义渐变色（可以自定义一组你喜欢的色号）
		colors := []string{"52", "88", "124", "125", "161", "198", "200", "201", "207"}

		// 渐变着色 Logo
		banner := gradient(logo, colors)

		yiyanstring := utils.Getyiyn()

		// 描述信息
		desc := lipgloss.NewStyle().
			Foreground(lipgloss.Color(utils.GetRodmoInt())). // 绿色
			Bold(true).
			Render("lain-cli 是一个基于 bubbletea + lipgloss 打造的现代 TUI 工具\n\n\t" + yiyanstring)

		return banner + "\n\n" + desc
	}(),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
