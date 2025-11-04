package exec

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	mui "github.com/rescheni/lain-cli/internal/ui"
	"github.com/rescheni/lain-cli/logs"
)

func Curl(method string, url string) {

	var respBody string
	var err error
	method = strings.ToUpper(method)
	// 当使用Post|PUT|PATCH的时候 输入json请求参数
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
		// get 等请求直接发送
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
}
