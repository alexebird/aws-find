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
}

//type SubcmdsConfig struct {
//Ec2 Ec2Config
//Ecr EcrConfig
//}

//type EcrConfig struct {
//Tableme TablemeConfig
//}

//type Ec2Config struct {
//Tableme TablemeConfig
//}

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
