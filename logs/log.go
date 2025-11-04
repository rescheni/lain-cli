package logs

import "fmt"

func Err(text string, err ...error) {
	fmt.Println("[ERROR]", text, err)
}

func Fatal(text string, err ...error) {

	panic(fmt.Sprintln("[FATAL]", text, err))
}

func Info(text string) {
	fmt.Println("[INFO]", text)
}

func Debug(text string) {
	fmt.Println("[DEBUG]", text)
}
