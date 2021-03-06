package ec2

import (
	"fmt"
	env "github.com/alexebird/aws-find/env"
	//"github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cliAction(c *cli.Context) error {
	client, err := setup()
	if err != nil {
		return err
	}

	instances := describeInstances(client, c.Bool("all"), c.String("filter"))
	noHeaders := c.GlobalBool("no-headers")

	if c.GlobalBool("color") == true {
		color.NoColor = false
	} else if c.GlobalBool("monochrome") == true {
		color.NoColor = true
	}

	//if c.Bool("all") == false && env.DavinciEnvFull() != nil {
	//davinciShortFormTable(instances, noHeaders)
	//}
	//else if c.Bool("all") == true && env.DavinciEnvFull() != nil {
	//davinciLongFormTable(instances, noHeaders)
	//}
	if c.Bool("all") == false && env.DavinciEnvFull() == nil {
		shortFormTable(instances, noHeaders)
	} else if c.Bool("all") == true && env.DavinciEnvFull() == nil {
		longFormTable(instances, noHeaders)
	}

	if c.Bool("connect") {
		inst := chooseInstanceForConnect(instances)
		if inst != nil {
			connect(inst)
		} else {
			fmt.Println("no connectable instances")
		}
	}

	return nil
}

func CliCommand() cli.Command {
	flags := []cli.Flag{
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "dont do default filtering",
		},
		cli.BoolFlag{
			Name:  "connect, c",
			Usage: "connect to the first matching host",
		},
		cli.StringFlag{
			Name:  "filter, f",
			Value: "",
			Usage: "filter by instance name",
		},
	}

	return cli.Command{
		Name:                   "ec2",
		Action:                 cliAction,
		Flags:                  flags,
		UseShortOptionHandling: true,
	}
}
