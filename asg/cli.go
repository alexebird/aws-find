package asg

import (
	//"fmt"

	//"github.com/davecgh/go-spew/spew"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/urfave/cli"
)

func cliAction(c *cli.Context) error {
	client, err := setup()
	if err != nil {
		return err
	}

	chanAsg := make(chan []*autoscaling.Group)
	chanLc := make(chan []*autoscaling.LaunchConfiguration)

	go func(c chan []*autoscaling.Group) {
		c <- describeAutoscalingGroups(client)
	}(chanAsg)

	go func(c chan []*autoscaling.LaunchConfiguration) {
		c <- describeLaunchConfigurations(client)
	}(chanLc)

	asgs := <-chanAsg
	lcs := <-chanLc
	//spew.Dump(asgs, lcs)

	printAutoScalingGroupsTable(asgs, lcs)

	//repo := c.String("repo")
	//all := c.Bool("all")
	//minTags := c.Int("min-tags")

	//if repo == "" {
	//repos := describeRepositories(client)
	//printReposTable(repos)
	//} else {
	//images := describeImages(client, repo, all)
	//printImagesTable(images, minTags)
	//}

	return nil
}

func CliCommand() cli.Command {
	//flags := []cli.Flag{
	//cli.BoolFlag{
	//Name:  "all, a",
	//Usage: "show all images",
	//},
	//cli.StringFlag{
	//Name:  "repo, r",
	//Usage: "show images for repo",
	//},
	//cli.IntFlag{
	//Name:  "min-tags, m",
	//Value: 0,
	//Usage: "only show images with m or more tags",
	//},
	//}

	return cli.Command{
		Name:   "asg",
		Action: cliAction,
		//Flags:  flags,
		UseShortOptionHandling: true,
	}
}
