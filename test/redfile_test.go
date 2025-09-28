package test

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestReadfile(t *testing.T) {

	wd, _ := os.Getwd()
	fmt.Println("当前工作目录:", wd)

	file, err := os.Open("../prompt/system_prompt.md")
	if err != nil {
		t.Fatal("open file Error", err)
	}
	defer file.Close()
	buff := make([]byte, 10)

	var prompt []byte

	for {
		n, err := file.Read(buff)
		if err != nil {
			if err != io.EOF {
				t.Fatal("read file error ")
			}
			break
		}
		prompt = append(prompt, buff[:n]...)

	}

	fmt.Println(string(prompt))

}
