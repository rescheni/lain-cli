package mui

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// responseMsg 是一个 "消息", 在 HTTP 请求完成后发回
type responseMsg struct {
	status int
	body   string
	err    error
}

// model_input 是我们 TUI 应用的 "状态"
type model_input struct {
	url      string         // 目标 URL
	method   string         // HTTP 方法 (POST, PUT, PATCH)
	textarea textarea.Model // "窗口" (JSON 编辑器)
	sending  bool           // 是否正在发送...
	response string         // 保存服务器的响应
	err      error          // 保存错误
	aborted  bool           // 用户是否按 Ctrl+C 中止了
}

// initialModel 创建我们的初始状态
// (已修改: 接收 method)
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

			// (已修改: 传递 method)
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

		// 成功！保存响应
		m.response = fmt.Sprintf("HTTP Status: %d\n\n%s", msg.status, msg.body)

		// (已修改) 成功后，我们也直接退出 TUI
		// TUI 会在退出时把 response 返回给调用者
		return m, tea.Quit

	// 消息类型：窗口大小调整
	case tea.WindowSizeMsg:
		m.textarea.SetWidth(msg.Width - 2)
	}

	return m, tea.Batch(cmd, taCmd)
}

// --- Bubble Tea View (渲染 UI) ---

var (
	titleStyle = lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230")).Padding(0, 1)
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Padding(1, 0)
	respStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
)

func (m model_input) View() string {
	// (已修改) View 不再需要显示 response, 因为成功后 TUI 会立刻退出

	// 2. 如果正在发送...
	if m.sending {
		return "\n  Sending request to " + m.url + " ..."
	}

	// 3. 默认视图 (编辑器)
	var b strings.Builder

	// (已修改: 标题显示 Method)
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

// sendRequestCmd (已修改: 接收 method)
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

// --- 痛点 3 & 4 修复: 核心入口函数 ---

// OpenTextView (已修改)
// 这是你的 lain-cli 应该调用的唯一函数
// 它返回 (输入的 JSON, 返回的 Body, 错误)
func OpenTextView(method string, url string) (string, string, error) {

	p := tea.NewProgram(initialModel(url, method), tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Println("启动 TUI 失败:", err)
		return "", "", err
	}
	m, ok := finalModel.(model_input)
	if !ok {
		return "", "", errors.New("TUI 模型转换失败")
	}

	// 检查是否是用户主动中止
	if m.aborted {
		return "", "", errors.New("请求中止")
	}

	// 检查是否成功（即 response 是否被填入）
	if m.response == "" {
		// 这可能发生在 JSON 校验失败时用户退出了 (虽然我们标记了aborted, 但这是双重保险)
		return "", "", errors.New("未收到响应")
	}

	// 成功！返回 TUI 里的数据
	// m.textarea.Value() 是你输入的历史 JSON
	// m.response 是服务器返回的结果
	return m.textarea.Value(), m.response, nil
}
