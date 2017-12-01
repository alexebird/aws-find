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

	instances := describeInstances(client, c.Bool("all"))

	if c.Bool("all") == false && d.DavinciEnvFull() != nil {
		davinciShortFormTable(instances)
	} else if c.Bool("all") == false && d.DavinciEnvFull() == nil {
		shortFormTable(instances)
	} else if c.Bool("all") == true && d.DavinciEnvFull() != nil {
		davinciLongFormTable(instances)
	} else if c.Bool("all") == true && d.DavinciEnvFull() == nil {
		longFormTable(instances)
	}

	return nil
}

func CliCommand() cli.Command {
	//var all bool
	flags := []cli.Flag{
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "dont do default filtering",
			//Destination: &all,
		},
	}

	return cli.Command{
		Name:   "ec2",
		Action: cliAction,
		Flags:  flags,
	}
}
