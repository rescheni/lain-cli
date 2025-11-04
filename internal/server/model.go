package server

import (
	"context"
	"errors"
	"fmt"

	config "github.com/rescheni/lain-cli/config"
	"github.com/rescheni/lain-cli/internal/tools"
	mui "github.com/rescheni/lain-cli/internal/ui"

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

func CallModel(ctx context.Context, ask string, useWindwos bool) error {

	completion, err := llms.GenerateFromSinglePrompt(ctx,
		LLLM,
		ask,
		llms.WithTemperature(0.8),
		llms.WithStopWords([]string{"Armstrong"}),
	)
	if err != nil {
		return errors.New("server CallModel error")
	}

	if config.Conf.Context.Enabled {
		tools.LLMCTX.Add(completion)
	}

	mui.PrintMarkdown(completion, useWindwos)
	return nil
}

func CallModelStream(ctx context.Context, ask string) (err error) {

	ai_ctx := ""
	_, err = LLLM.Call(
		ctx,
		ask,
		llms.WithTemperature(0.7),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			schunk := string(chunk)
			fmt.Print(schunk) // 每次收到一段 token 就打印
			if config.Conf.Context.Enabled {
				ai_ctx += schunk
			}
			return nil
		}),
	)
	if err != nil {
		return err
	}
	if config.Conf.Context.Enabled {
		// 写入上下文
		tools.LLMCTX.Add(ai_ctx)
	}

	return nil
}
