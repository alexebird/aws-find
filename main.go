package main

import (
	//"fmt"
	"os"

	afec2 "github.com/alexebird/aws-find/ec2"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "aws-find"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		afec2.CliCommand(),
	}

	app.Run(os.Args)
}
