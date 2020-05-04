package ami

import (
	//"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cliAction(c *cli.Context) error {
	client, err := setup()
	if err != nil {
		return err
	}

	images := describeImages(client)
	noHeaders := c.GlobalBool("no-headers")

	//spew.Dump(images)

	if c.GlobalBool("color") == true {
		color.NoColor = false
	} else if c.GlobalBool("monochrome") == true {
		color.NoColor = true
	}

	//if c.Bool("all") == false && env.DavinciEnvFull() == nil {
	shortFormTable(images, noHeaders)
	//}
	//} else if c.Bool("all") == true && env.DavinciEnvFull() == nil {
	//longFormTable(instances, noHeaders)
	//}

	return nil
}

func CliCommand() cli.Command {
	//flags := []cli.Flag{
	//cli.BoolFlag{
	//Name:  "all, a",
	//Usage: "dont do default filtering",
	//},
	//cli.BoolFlag{
	//Name:  "connect, c",
	//Usage: "connect to the first matching host",
	//},
	//cli.StringFlag{
	//Name:  "filter, f",
	//Value: "",
	//Usage: "filter by instance name",
	//},
	//}

	return cli.Command{
		Name:   "ami",
		Action: cliAction,
		//Flags:                  flags,
		UseShortOptionHandling: true,
	}
}
