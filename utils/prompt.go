package utils

import (
	"lain-cli/logs"
	"os"
)

func Getprompt() string {

	data, err := os.ReadFile("prompt/system_prompt.md")
	if err != nil {
		logs.Err("open os prompt error:", err)
		return ""
	}
	return string(data)
}
