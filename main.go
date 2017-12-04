package main

import (
	//"fmt"
	"os"

	afec2 "github.com/alexebird/aws-find/ec2"
	afecr "github.com/alexebird/aws-find/ecr"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "aws-find"
	app.Version = "0.0.2"

	app.Commands = []cli.Command{
		afec2.CliCommand(),
		afecr.CliCommand(),
	}

	app.Run(os.Args)
}
