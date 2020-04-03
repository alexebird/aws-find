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

	noHeaders := c.GlobalBool("no-headers")

	if repo == "" {
		repos := describeRepositories(client)
		printReposTable(repos, noHeaders)
	} else {
		images := describeImages(client, repo, all)
		printImagesTable(images, minTags, noHeaders)
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
		Name:                   "ecr",
		Action:                 cliAction,
		Flags:                  flags,
		UseShortOptionHandling: true,
	}
}
