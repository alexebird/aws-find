package main

import (
	//"fmt"
	"os"
	"path"

	afasg "github.com/alexebird/aws-find/asg"
	config "github.com/alexebird/aws-find/config"
	afec2 "github.com/alexebird/aws-find/ec2"
	afecr "github.com/alexebird/aws-find/ecr"
	//"github.com/davecgh/go-spew/spew"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "aws-find"
	app.Version = "0.0.2"

	flags := []cli.Flag{
		cli.BoolFlag{
			Name:  "color, C",
			Usage: "force color",
		},
		cli.BoolFlag{
			Name:  "monochrome, M",
			Usage: "force no color",
		},
		cli.BoolFlag{
			Name:  "no-headers, H",
			Usage: "dont print headers",
		},
	}
	app.Flags = flags

	app.Commands = []cli.Command{
		afec2.CliCommand(),
		afecr.CliCommand(),
		afasg.CliCommand(),
	}

	home, _ := homedir.Dir()
	confFile := path.Join(home, ".aws-find.yml")

	if _, err := os.Stat(confFile); err == nil {
		config := config.ReadConfig(confFile)
		afasg.Config = config
		afec2.Config = config
		afecr.Config = config
	}

	app.Run(os.Args)
}
