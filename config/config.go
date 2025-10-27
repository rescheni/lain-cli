package config

import (
	"os"

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
type Config struct {
	Ai      Ai      `yaml:"ai"`
	Yiyan   Yiyan   `yaml:"yiyan"`
	Context Context `yaml:"context"`
	Mcp     Mcp     `yaml:"mcp"`
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
	// fmt.Println(Conf)
}
