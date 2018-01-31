package config

import (
	//"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var Config AwsFindConfig

type AwsFindConfig struct {
	Tableme TablemeConfig
	Ec2     Ec2Config
}

//type SubcmdsConfig struct {
//Ec2 Ec2Config
//Ecr EcrConfig
//}

//type EcrConfig struct {
//Tableme TablemeConfig
//}

type AutofilterConfig struct {
	Tag    string
	Values []EnvVarConfig
}

//type ValuesConfig struct {
//EnvVars []EnvVarConfig
//}

type EnvVarConfig struct {
	EnvVar string `yaml:"env_var"`
}

type Ec2Config struct {
	Autofilter AutofilterConfig
}

type TablemeConfig struct {
	Colorize []ColorizeConfig
}

type ColorizeConfig struct {
	Subcmds []string
	Regex   string
	Color   string
}

func ReadConfig(confPath string) AwsFindConfig {
	config := AwsFindConfig{}

	content, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return config
}
