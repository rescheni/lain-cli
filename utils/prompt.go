package utils

import (
	"fmt"
	"os"
)

func Getprompt() string {

	data, err := os.ReadFile("prompt/system_prompt.md")
	if err != nil {
		fmt.Println("open os prompt error:", err)
		return ""
	}
	return string(data)
}
