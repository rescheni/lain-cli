package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	spin     spinner.Model
	quitting bool
}

func (m model) Init() tea.Cmd {
	return m.spin.Tick // 启动定时器
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spin, cmd = m.spin.Update(msg) // 更新 spinner 状态
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "✔ 完成!\n"
	}
	return fmt.Sprintf("%s  正在思考...\n", m.spin.View())
}

func TestThinkui(T *testing.T) {
	spin := spinner.New()
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	spin.Spinner = spinner.Dot

	m := model{spin: spin}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
