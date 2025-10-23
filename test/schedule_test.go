package test

import (
	"fmt"
	"lain-cli/exec"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 0
	maxWidth = 100
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

func TestSchedule(t *testing.T) {
	m := models{
		progress: progress.New(progress.WithDefaultGradient()),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
	fmt.Println()
	exec.OutScanner()
}

func runScanCmd() tea.Cmd {
	return func() tea.Msg {
		exec.RunNmap("reschen.cn", 20, 10000)
		return "scan finished" // (你可以自定义一个 struct)
	}
}

type tickMsg time.Time

type models struct {
	progress progress.Model
}

func (m models) Init() tea.Cmd {

	return tea.Batch(
		runScanCmd(),
		tickCmd(),
	)
}

func (m models) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:

		currentDone := exec.CompletedTasks.Load()
		percent := float64(currentDone) / float64(exec.TotalTasks)
		if percent >= 1.0 {

			return m, tea.Quit
		}
		cmd := m.progress.SetPercent(percent)
		return m, tea.Batch(tickCmd(), cmd)

	case progress.FrameMsg:
		progressmodels, cmd := m.progress.Update(msg)
		m.progress = progressmodels.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m models) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + m.progress.View() + "\n\n" +
		pad + helpStyle("Press ctrl + c to quit")
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
