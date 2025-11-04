package utils

import (
	"os"

	logs "github.com/rescheni/lain-cli/logs"
)

func Getprompt() string {

	data, err := os.ReadFile("prompt/system_prompt.md")
	if err != nil {
		logs.Err("open os prompt error:", err)
		return ""
	}
	return string(data)
}
