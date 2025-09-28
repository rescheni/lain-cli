package utils

import (
	"encoding/json"
	"io"
	"lain-cli/config"
	"net/http"
)

type info struct {
	Hitokoto string `json:"hitokoto"`
	From     string `json:"from"`
}

func Getyiyn() string {

	status := config.Conf.Yiyan.Status
	url := config.Conf.Yiyan.Api_url
	if status != "ON" {
		return ""
	}

	var infoo info
	resp, err := http.Get(url + "?encode=json")
	if err != nil {
		return "一言接口错误/网络错误"
	}
	defer resp.Body.Close()
	s, err := io.ReadAll(resp.Body)
	if err != nil {
		return "parse body err"
	}

	err = json.Unmarshal(s, &infoo)
	if err != nil {
		return "json.Unmarshal err "
	}

	return infoo.Hitokoto + "\n\t\t\t\t\t\t————" + infoo.From
}
