package ami

import (
	"fmt"
	"sort"
	"time"

	config "github.com/alexebird/aws-find/config"
	env "github.com/alexebird/aws-find/env"
	util "github.com/alexebird/aws-find/util"
	"github.com/alexebird/tableme/tableme"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	//"github.com/davecgh/go-spew/spew"
)

var Config config.AwsFindConfig

const longForm = "2006-01-02T15:04:05.999Z"

type ByCreationDate []*ec2.Image

func (s ByCreationDate) Len() int {
	return len(s)
}
func (s ByCreationDate) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByCreationDate) Less(i, j int) bool {
	//CreationDate: "2020-05-01T21:08:23.000Z",

	ti := *s[i].CreationDate
	//spew.Dump(ti)
	ti_p, _ := time.Parse(longForm, ti)
	//spew.Dump(ti_p)

	tj := *s[j].CreationDate
	tj_p, _ := time.Parse(longForm, tj)

	// reverse sort
	return ti_p.Before(tj_p)
}

func describeImages(client *ec2.EC2) []*ec2.Image {
	svc := ec2.New(session.New())
	currentOwnerId := env.GetEnv("AWS_ACCOUNT_ID")
	input := &ec2.DescribeImagesInput{
		Owners: []*string{
			currentOwnerId,
		},
	}

	result, err := svc.DescribeImages(input)
	if err != nil {
		panic(err)
	}

	sort.Sort(ByCreationDate(result.Images))

	return result.Images

	//var filters []*ec2.Filter

	//if all == true {
	//filters = nil
	//} else {
	//filters = configFilters()
	//}

	//params := &ec2.DescribeInstancesInput{Filters: filters}
	//instances := make([]*ec2.Instance, 0)

	//err := client.DescribeInstancesPages(params,
	//func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
	//for _, res := range page.Reservations {
	//instances = append(instances, res.Instances...)
	//}
	//return true
	//})

	//if err != nil {
	//panic(err)
	//}

	//if len(nameFilter) > 0 {
	//filterTest := func(i *ec2.Instance) bool {
	//nameTagValue := findTagByKey(i, "Name")
	//if nameTagValue != nil {
	//return strings.Contains(strings.ToLower(*nameTagValue), strings.ToLower(nameFilter))
	//} else {
	//return false
	//}
	//}
	//instances = filterInstances(instances, filterTest)
	//}

	//sort.Sort(ByCreationDate(instances))

	//return instances
}

func shortFormTable(images []*ec2.Image, noHeaders bool) {
	headers := []string{
		"NAME", "STATE", "PUBLIC", "OWNER_ID", "CREATION", "IMAGE_ID",
	}

	records := make([][]string, 0)

	for _, img := range images {
		t, _ := time.Parse(longForm, *img.CreationDate)
		elapsed := time.Now().Sub(t).Hours() / 24

		rec := []string{
			tableme.WithEmptyStringDefault(img.Name),
			tableme.WithEmptyStringDefault(img.State),
			tableme.StringifyBool(*img.Public),
			tableme.WithEmptyStringDefault(img.OwnerId),
			tableme.StringifyString(fmt.Sprintf("%3.0fd %s", elapsed, *img.CreationDate)),
			tableme.WithEmptyStringDefault(img.ImageId),
		}
		records = append(records, rec)
	}

	bites := tableme.TableMe(headers, records, noHeaders)
	util.PrintColorizedTable(bites, "ami", Config.Tableme.Colorize)
}

func setup() (*ec2.EC2, error) {
	var client *ec2.EC2 = ec2.New(session.New())
	return client, nil
}
