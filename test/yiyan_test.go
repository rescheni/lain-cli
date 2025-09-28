package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

type info struct {
	Hitokoto string `json:"hitokoto"`
	From     string `json:"from"`
}

func TestGetyiyn(t *testing.T) {
	var infoo info

	resp, err := http.Get("https://v1.hitokoto.cn?encode=json")
	if err != nil {
		t.Error("yiayan error")
		return
	}
	defer resp.Body.Close()
	s, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error("parse body err")
		return
	}

	err = json.Unmarshal(s, &infoo)
	if err != nil {
		t.Error("json.Unmarshal err ")
		return
	}

	fmt.Println(infoo.Hitokoto + "\n\t\t\t————" + infoo.From)
}
