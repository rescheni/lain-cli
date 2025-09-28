package server

import (
	"context"
	"errors"
	"fmt"
	"lain-cli/config"

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

func CallModel(ctx context.Context, ask string) (string, error) {

	completion, err := llms.GenerateFromSinglePrompt(ctx,
		LLLM,
		ask,
		llms.WithTemperature(0.8),
		llms.WithStopWords([]string{"Armstrong"}),
	)
	if err != nil {
		return "", errors.New("server CallModel error")
	}
	return completion, nil
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
