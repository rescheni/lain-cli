package server

import (
	"context"
	"errors"
	"fmt"
	"lain-cli/config"
	mui "lain-cli/ui"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

var LLLM *openai.LLM

func LLMInit() error {

	// apikey := os.Getenv("OPENAI_API_KEY")
	// apiadd := os.Getenv("OPENAI_API_ADD")
	// aimodelname := os.Getenv("AI_MODEL_NAME")
	apikey := config.Conf.Ai.Api_key
	apiadd := config.Conf.Ai.Api_url
	aimodelname := config.Conf.Ai.Api_model_name

	llm, err := openai.New(
		openai.WithModel(aimodelname),
		openai.WithToken(apikey),
		openai.WithBaseURL(apiadd),
	)
	if err != nil {
		return errors.New("init AI model err")
	}
	LLLM = llm
	return nil

}

func CallModel(ctx context.Context, ask string) error {

	completion, err := llms.GenerateFromSinglePrompt(ctx,
		LLLM,
		ask,
		llms.WithTemperature(0.8),
		llms.WithStopWords([]string{"Armstrong"}),
	)
	if err != nil {
		return errors.New("server CallModel error")
	}

	// 我们把要渲染的 Markdown 内容传进去
	m := mui.NewMarkdownModel(completion)

	// 2. 创建一个新的 Bubble Tea 程序
	// tea.WithAltScreen() 进入“全屏”模式
	// tea.WithMouseCellMotion() 启用鼠标滚轮支持
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// 3. 运行！
	// .Run() 会接管终端，直到用户按 Ctrl+C
	if _, err := p.Run(); err != nil {
		fmt.Println("启动 TUI 失败:", err)
		os.Exit(1)
	}
	return nil
}

func CallModelStream(ctx context.Context, ask string) (err error) {

	_, err = LLLM.Call(
		ctx,
		ask,
		llms.WithTemperature(0.7),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk)) // 每次收到一段 token 就打印
			return nil
		}),
	)
	if err != nil {
		return err
	}

	return nil
}
