package config

import (
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Ai struct {
	Api_url        string `yaml:"api_url"`
	Api_key        string `yaml:"api_key"`
	Api_model_name string `yaml:"ai_model_name"`
}

type Yiyan struct {
	Status  string `yaml:"status"`
	Api_url string `yaml:"api_url"`
}

type Context struct {
	Max_length int    `yaml:"max_length"`
	Local      string `yaml:"local"`
	Enabled    bool   `yaml:"enabled"`
}

type Mcp struct {
	Json string `yaml:"json"`
}
type Logo struct {
	Logo_txt string `yaml:"logo_txt"`
}
type Config struct {
	Ai      Ai      `yaml:"ai"`
	Yiyan   Yiyan   `yaml:"yiyan"`
	Context Context `yaml:"context"`
	Mcp     Mcp     `yaml:"mcp"`
	Logo    Logo    `yaml:"logo"`
}

var Conf Config

func init() {
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		return
	}

	err = yaml.Unmarshal(data, &Conf)
	if err != nil {
		return
	}
	Check_ENV(&Conf.Ai.Api_key)
	Check_ENV(&Conf.Ai.Api_model_name)
	Check_ENV(&Conf.Ai.Api_url)
	Check_ENV(&Conf.Context.Local)

	Check_ENV(&Conf.Mcp.Json)
	Check_ENV(&Conf.Yiyan.Api_url)

}

func Check_ENV(conf *string) {
	if strings.HasPrefix(*conf, "ENV_") {
		*conf = os.Getenv(*conf)
	}
}
