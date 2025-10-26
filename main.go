/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"lain-cli/cmd"
	"lain-cli/server"
)

func main() {

	err := server.LLMInit()
	if err != nil {
		fmt.Println("server LLM init Error")
	}
	cmd.Execute()
}
