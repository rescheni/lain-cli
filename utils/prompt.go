package utils

import (
	"fmt"
	"os"
)

var Prompt string

func Initprompt() {

	data, err := os.ReadFile("prompt/system_prompt.md")
	if err != nil {
		fmt.Println("open os prompt error:", err)
		return
	}
	Prompt = string(data)
}
