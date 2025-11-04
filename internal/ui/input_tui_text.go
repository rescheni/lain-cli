package mui

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/rescheni/lain-cli/logs"

	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// responseMsg 在 HTTP 请求完成后发回
type responseMsg struct {
	status int
	body   string
	err    error
}

// TUI 应用状态
type model_input struct {
	url      string         // 目标 URL
	method   string         // HTTP 方法 (POST, PUT, PATCH)
	textarea textarea.Model // "窗口" (JSON 编辑器)
	sending  bool           // 是否正在发送...
	response string         // 保存服务器的响应
	err      error          // 保存错误
	aborted  bool           // 用户是否按 Ctrl+C 中止了
}

// initialModel 创建初始状态
func initialModel(url string, method string) model_input {
	ta := textarea.New()
	ta.Placeholder = "Place Input json\n{\n  \"key\": \"value\"\n}\n..."
	ta.Focus()
	ta.SetHeight(10)
	ta.SetWidth(50)
	ta.ShowLineNumbers = true
	ta.CharLimit = 0

	return model_input{
		url:      url,
		method:   method,
		textarea: ta,
		sending:  false,
		err:      nil,
		aborted:  false,
	}
}
func (m model_input) Init() tea.Cmd {
	return textarea.Blink
}

func (m model_input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		taCmd tea.Cmd
		cmd   tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			m.textarea.InsertString("	")
			return m, nil
		}
	}

	m.textarea, taCmd = m.textarea.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.aborted = true
			return m, tea.Quit

		case tea.KeyCtrlS:
			if m.sending {
				return m, nil
			}
			jsonStr := m.textarea.Value()

			if !json.Valid([]byte(jsonStr)) {
				m.err = errors.New("JSON 格式无效。请修正后再发送。")
				return m, nil
			}

			m.err = nil
			m.sending = true
			m.textarea.Blur()

			return m, sendRequestCmd(m.url, m.method, jsonStr)

		default:
			m.err = nil
		}

	// 消息类型：我们的自定义 HTTP 响应
	case responseMsg:
		m.sending = false

		if msg.err != nil {
			m.err = msg.err
			m.textarea.Focus()
			return m, nil
		}

		m.response = fmt.Sprintf("HTTP Status: %d\n\n%s", msg.status, msg.body)
		return m, tea.Quit

	// 窗口大小调整
	case tea.WindowSizeMsg:
		m.textarea.SetWidth(msg.Width - 2)
	}

	return m, tea.Batch(cmd, taCmd)
}

// --- Bubble Tea View (渲染 UI) ---

var (
	titleStyle = lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230")).Padding(0, 1)
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Padding(1, 0)
	// respStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
)

func (m model_input) View() string {

	if m.sending {
		return "\n  Sending request to " + m.url + " ..."
	}

	var b strings.Builder

	b.WriteString(titleStyle.Render(fmt.Sprintf("%s to %s", m.method, m.url)))
	b.WriteString("\n\n")
	b.WriteString(m.textarea.View())
	b.WriteString("\n\n")
	b.WriteString("编辑 JSON (按 Esc/Ctrl+C 退出, 按 Ctrl+S 发送, 按 Tab 插入空格)")

	if m.err != nil {
		b.WriteString("\n" + errorStyle.Render(m.err.Error()))
	}

	return b.String()
}

func sendRequestCmd(url, method, jsonStr string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer([]byte(jsonStr)))
		if err != nil {
			return responseMsg{err: err}
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return responseMsg{err: err}
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return responseMsg{err: err}
		}

		return responseMsg{
			status: resp.StatusCode,
			body:   string(body),
			err:    nil,
		}
	}
}

func OpenTextView(method string, url string) (string, string, error) {

	p := tea.NewProgram(initialModel(url, method), tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		logs.Err("启动 TUI 失败:", err)
		return "", "", err
	}
	m, ok := finalModel.(model_input)
	if !ok {
		return "", "", errors.New("TUI 模型转换失败")
	}

	if m.aborted {
		return "", "", errors.New("请求中止")
	}

	if m.response == "" {
		return "", "", errors.New("未收到响应")
	}

	return m.textarea.Value(), m.response, nil
}
