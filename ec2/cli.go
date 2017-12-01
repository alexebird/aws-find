package ec2

import (
	d "github.com/alexebird/aws-find/davinci"
	//"github.com/davecgh/go-spew/spew"
	"github.com/urfave/cli"
)

func cliAction(c *cli.Context) error {
	//spew.Dump(c.Bool("all"))
	client, err := setup()
	if err != nil {
		return err
	}

	instances := describeInstances(client, c.Bool("all"), c.String("filter"))

	if c.Bool("all") == false && d.DavinciEnvFull() != nil {
		davinciShortFormTable(instances)
	} else if c.Bool("all") == false && d.DavinciEnvFull() == nil {
		shortFormTable(instances)
	} else if c.Bool("all") == true && d.DavinciEnvFull() != nil {
		davinciLongFormTable(instances)
	} else if c.Bool("all") == true && d.DavinciEnvFull() == nil {
		longFormTable(instances)
	}

	if c.Bool("connect") {
		connect(instances)
	}

	return nil
}

func CliCommand() cli.Command {
	//var all bool
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
		Name:   "ec2",
		Action: cliAction,
		Flags:  flags,
	}
}
