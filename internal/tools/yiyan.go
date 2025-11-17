package tools

import (
	"encoding/json"
	"io"
	"net/http"

	config "github.com/rescheni/lain-cli/config"
)

type yy_info struct {
	Hitokoto string `json:"hitokoto"`
	From     string `json:"from"`
}

func Getyiyn() string {

	// 检查配置文件 一言开关是否打开？
	status := config.Conf.Yiyan.Status
	url := config.Conf.Yiyan.Api_url
	if status != "ON" {
		return ""
	}
	// 发送get请求 json 数据
	var infoo yy_info
	resp, err := http.Get(url + "?encode=json")
	if err != nil {
		return "一言接口错误/网络错误"
	}
	defer resp.Body.Close()
	s, err := io.ReadAll(resp.Body)
	if err != nil {
		return "parse body err"
	}
	// 	解析数据
	err = json.Unmarshal(s, &infoo)
	if err != nil {
		return "json.Unmarshal err "
	}
	return infoo.Hitokoto + "\n\t\t\t\t\t\t————" + infoo.From
}
