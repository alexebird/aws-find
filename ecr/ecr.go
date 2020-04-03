package ecr

import (
	"fmt"
	"sort"
	"strings"
	"time"

	config "github.com/alexebird/aws-find/config"
	env "github.com/alexebird/aws-find/env"
	util "github.com/alexebird/aws-find/util"
	"github.com/alexebird/tableme/tableme"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	//"github.com/davecgh/go-spew/spew"
)

var Config config.AwsFindConfig

type ByRepositoryName []*ecr.Repository

func (s ByRepositoryName) Len() int {
	return len(s)
}
func (s ByRepositoryName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByRepositoryName) Less(i, j int) bool {
	return strings.Compare(*s[i].RepositoryName, *s[j].RepositoryName) < 0
}

type ByImagePushedAt []*ecr.ImageDetail

func (s ByImagePushedAt) Len() int {
	return len(s)
}
func (s ByImagePushedAt) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByImagePushedAt) Less(i, j int) bool {
	ti := *s[i].ImagePushedAt
	tj := *s[j].ImagePushedAt
	return tj.After(ti)
}

func describeRepositories(client *ecr.ECR) []*ecr.Repository {
	params := &ecr.DescribeRepositoriesInput{}
	repos := make([]*ecr.Repository, 0)

	err := client.DescribeRepositoriesPages(params,
		func(page *ecr.DescribeRepositoriesOutput, lastPage bool) bool {
			repos = append(repos, page.Repositories...)
			return true
		})

	if err != nil {
		panic(err)
	}

	return repos
}

func printReposTable(repos []*ecr.Repository, noHeaders bool) {
	sort.Sort(ByRepositoryName(repos))

	headers := []string{
		"NAME",
	}

	records := make([][]string, 0)

	for _, inst := range repos {
		rec := []string{
			tableme.WithEmptyStringDefault(inst.RepositoryName),
		}
		records = append(records, rec)
	}

	bites := tableme.TableMe(headers, records, noHeaders)
	util.PrintColorizedTable(bites, "ecr", Config.Tableme.Colorize)
}

func describeImages(client *ecr.ECR, repo string, all bool) []*ecr.ImageDetail {
	var filter *ecr.DescribeImagesFilter

	if all {
		filter = nil
	} else {
		filter = &ecr.DescribeImagesFilter{
			TagStatus: aws.String("TAGGED"),
		}
	}

	params := &ecr.DescribeImagesInput{
		RegistryId:     env.MustGetEnv("ECR_REGISTRY_ID"),
		RepositoryName: &repo,
		Filter:         filter,
	}

	images := make([]*ecr.ImageDetail, 0)

	err := client.DescribeImagesPages(params,
		func(page *ecr.DescribeImagesOutput, lastPage bool) bool {
			images = append(images, page.ImageDetails...)
			return true
		})

	if err != nil {
		panic(err)
	}

	return images
}

func imgTagString(tags []*string) string {
	newTags := make([]string, 0)
	for _, tag := range tags {
		newTags = append(newTags, *tag)
	}
	sort.Strings(newTags)
	return strings.Join(newTags, ",")
}

func printImagesTable(repos []*ecr.ImageDetail, minTags int, noHeaders bool) {
	sort.Sort(ByImagePushedAt(repos))

	headers := []string{
		"TAGS", "PUSHED", "SIZE", "SHA256",
	}

	records := make([][]string, 0)

	for _, img := range repos {
		if len(img.ImageTags) < minTags {
			continue
		}

		pushedAt := img.ImagePushedAt.Format(time.RFC3339)
		imgSize := fmt.Sprintf("%.2fMB", float64(*img.ImageSizeInBytes)/1024.0/1024.0)
		imgTags := imgTagString(img.ImageTags)

		rec := []string{
			tableme.WithEmptyStringDefault(&imgTags),
			tableme.WithEmptyStringDefault(&pushedAt),
			tableme.WithEmptyStringDefault(&imgSize),
			tableme.WithEmptyStringDefault(img.ImageDigest),
		}
		records = append(records, rec)
	}

	bites := tableme.TableMe(headers, records, noHeaders)
	util.PrintColorizedTable(bites, "ecr", Config.Tableme.Colorize)
}

func setup() (*ecr.ECR, error) {
	var client *ecr.ECR = ecr.New(session.New(), aws.NewConfig().WithRegion(*env.MustGetEnv("ECR_REGION")))
	return client, nil
}
