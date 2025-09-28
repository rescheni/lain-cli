package test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func TestModel(t *testing.T) {
	apiadd := "https://ai.reschen.cn:888/api"
	apikey := "sk-5b0cc6383a554883acedb671f1704733"

	llm, err := openai.New(
		openai.WithModel("qwen3:1.7b"),
		openai.WithToken(apikey),
		openai.WithBaseURL(apiadd),
	)
	if err != nil {
		t.Error("llm init error ")
	}

	ctx := context.Background()

	completion, err := llms.GenerateFromSinglePrompt(ctx,
		llm,
		"第一个在月球上行走的人",
		llms.WithTemperature(0.8),
		// llms.WithStopWords([]string{"Armstrong"}),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("The first man to walk on the moon:")
	fmt.Println(completion)

	t.Log(completion)
}
