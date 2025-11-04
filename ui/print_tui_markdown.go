package mui

import (
	"lain-cli/logs"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

var markdown string

type model_markdown struct {
	viewport        viewport.Model
	renderer        *glamour.TermRenderer
	markdownContent string
}

func PrintMarkdown(completion string, uw bool) {
	// 我们把要渲染的 Markdown 内容传进去
	m := NewMarkdownModel(completion)
	markdown = completion
	var p *tea.Program
	if uw {
		p = tea.NewProgram(
			m,
			tea.WithAltScreen(),
			tea.WithMouseCellMotion(),
		)
	} else {
		p = tea.NewProgram(
			m,
			// tea.WithAltScreen(),
			tea.WithMouseCellMotion(),
		)
	}

	// 3. 运行！
	// .Run() 会接管终端，直到用户按 Ctrl+C
	if _, err := p.Run(); err != nil {
		logs.Err("启动 TUI 失败:", err)
		os.Exit(1)
	}
}

func NewMarkdownModel(content string) model_markdown {
	// 1. 创建 viewport
	vp := viewport.New(50, 5) // 初始大小

	// 2. 创建 renderer
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(60), // 初始宽度
	)
	if err != nil {
		log.Fatal(err)
	}

	// 3. 渲染一次
	styledContent, _ := renderer.Render(content)
	vp.SetContent(styledContent)

	return model_markdown{
		viewport:        vp,
		renderer:        renderer,
		markdownContent: content,
	}

}

func (m model_markdown) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height

		renderer, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(msg.Width),
		)
		if err != nil {
			log.Println("failed to recreate glamour renderer:", err)
		} else {
			m.renderer = renderer
		}

		styledMarkdown, _ := m.renderer.Render(m.markdownContent)
		m.viewport.SetContent(styledMarkdown)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			return m, tea.Quit
		}
	}
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model_markdown) View() string {
	return m.viewport.View()
}

func (m model_markdown) Init() tea.Cmd {

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(50), // 初始宽度
	)
	styledMarkdown, _ := renderer.Render(markdown)
	m.viewport.SetContent(styledMarkdown)
	return nil
}
