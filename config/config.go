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

type Config struct {
	Ai    Ai    `yaml:"ai"`
	Yiyan Yiyan `yaml:"yiyan"`
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
