package exec

import (
	"fmt"
	"io"
	"lain-cli/logs"
	mui "lain-cli/ui"
	"net/http"
	"strings"
)

func Curl(method string, url string) error {

	var respBody string
	var err error
	method = strings.ToUpper(method)
	if method == "POST" || method == "PUT" || method == "PATCH" {
		var req string
		req, respBody, err = mui.OpenTextView(method, url)
		if err != nil {
			logs.Err("", err)
		}
		logs.Info("Your Ask:")
		fmt.Println(req)
		fmt.Println(respBody)
	} else {
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			logs.Err("", err)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logs.Err("http client:", err)
		}
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Println(string(bodyBytes))
	}
	return err
}
