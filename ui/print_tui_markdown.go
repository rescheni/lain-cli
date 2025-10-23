package mui

import (
	"log"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

var Markdown string

type model_markdown struct {
	viewport        viewport.Model
	renderer        *glamour.TermRenderer // 1. 存储渲染器
	markdownContent string                // 2. 存储原始 Markdown
}

// (你需要一个 New 函数来初始化)
func NewMarkdownModel(content string) model_markdown {
	// 1. 创建 viewport
	vp := viewport.New(80, 20) // 初始大小

	// 2. 创建 renderer
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80), // 初始宽度
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
		glamour.WithWordWrap(80), // 初始宽度
	)
	styledMarkdown, _ := renderer.Render(Markdown)
	m.viewport.SetContent(styledMarkdown)
	return nil
}
