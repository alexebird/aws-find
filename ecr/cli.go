package ecr

import (
	//"fmt"

	//"github.com/davecgh/go-spew/spew"
	"github.com/urfave/cli"
)

func cliAction(c *cli.Context) error {
	client, err := setup()
	if err != nil {
		return err
	}

	repo := c.String("repo")
	all := c.Bool("all")
	minTags := c.Int("min-tags")

	if repo == "" {
		repos := describeRepositories(client)
		printReposTable(repos)
	} else {
		images := describeImages(client, repo, all)
		printImagesTable(images, minTags)
	}

	return nil
}

func CliCommand() cli.Command {
	flags := []cli.Flag{
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "show all images",
		},
		cli.StringFlag{
			Name:  "repo, r",
			Usage: "show images for repo",
		},
		cli.IntFlag{
			Name:  "min-tags, m",
			Value: 0,
			Usage: "only show images with m or more tags",
		},
	}

	return cli.Command{
		Name:   "ecr",
		Action: cliAction,
		Flags:  flags,
		UseShortOptionHandling: true,
	}
}
