/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"lain-cli/cmd"
	"lain-cli/server"
	"lain-cli/utils"
)

func main() {

	err := server.LLMInit()
	if err != nil {
		fmt.Println("server LLM init Error")
	}
	utils.Initprompt()
	cmd.Execute()
}
