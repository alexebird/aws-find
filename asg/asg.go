package asg

import (
	"fmt"
	"sort"
	"strings"
	//"time"

	//env "github.com/alexebird/aws-find/env"
	config "github.com/alexebird/aws-find/config"
	util "github.com/alexebird/aws-find/util"
	"github.com/alexebird/tableme/tableme"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	//"github.com/davecgh/go-spew/spew"
)

var Config config.AwsFindConfig

type ByName []*autoscaling.Group

func (s ByName) Len() int {
	return len(s)
}
func (s ByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByName) Less(i, j int) bool {
	return strings.Compare(*s[i].AutoScalingGroupName, *s[j].AutoScalingGroupName) < 0
}

func describeAutoscalingGroups(client *autoscaling.AutoScaling) []*autoscaling.Group {
	params := &autoscaling.DescribeAutoScalingGroupsInput{}
	asgs := make([]*autoscaling.Group, 0)

	err := client.DescribeAutoScalingGroupsPages(params,
		func(page *autoscaling.DescribeAutoScalingGroupsOutput, lastPage bool) bool {
			asgs = append(asgs, page.AutoScalingGroups...)
			return true
		})

	if err != nil {
		panic(err)
	}

	return asgs
}

func describeLaunchConfigurations(client *autoscaling.AutoScaling) []*autoscaling.LaunchConfiguration {
	params := &autoscaling.DescribeLaunchConfigurationsInput{}
	lcs := make([]*autoscaling.LaunchConfiguration, 0)

	err := client.DescribeLaunchConfigurationsPages(params,
		func(page *autoscaling.DescribeLaunchConfigurationsOutput, lastPage bool) bool {
			lcs = append(lcs, page.LaunchConfigurations...)
			return true
		})

	if err != nil {
		panic(err)
	}

	return lcs
}

func mapLaunchConfigurationName(lcs []*autoscaling.LaunchConfiguration) map[string]*autoscaling.LaunchConfiguration {
	mapping := make(map[string]*autoscaling.LaunchConfiguration)

	for _, lc := range lcs {
		mapping[*lc.LaunchConfigurationName] = lc
	}

	return mapping
}

func printAutoScalingGroupsTable(asgs []*autoscaling.Group, lcs []*autoscaling.LaunchConfiguration, noHeaders bool) {
	sort.Sort(ByName(asgs))

	headers := []string{
		"NAME",
		"SIZE",
		"DESIRED",
		//"CREATED",
		//"LAUNCH CONFIG",
		"IMAGE",
		"TYPE",
	}

	mappedLcs := mapLaunchConfigurationName(lcs)
	records := make([][]string, 0)

	for _, asg := range asgs {
		//spew.Dump(mappedLcs)
		//spew.Dump(asg.LaunchConfigurationName)
		lc := mappedLcs[*asg.LaunchConfigurationName]
		//spew.Dump(lc)
		lcImage := lc.ImageId
		lcInstanceType := lc.InstanceType
		//time := asg.CreatedTime.Format(time.RFC3339)
		sizeStr := fmt.Sprintf("%d,%d", *asg.MinSize, *asg.MaxSize)

		rec := []string{
			tableme.StringifyStringPtr(asg.AutoScalingGroupName),
			tableme.StringifyString(sizeStr),
			tableme.StringifyIntPtr(asg.DesiredCapacity),
			//tableme.WithEmptyStringDefault(&time),
			//tableme.StringifyStringPtr(asg.LaunchConfigurationName),
			tableme.StringifyStringPtr(lcImage),
			tableme.StringifyStringPtr(lcInstanceType),
		}
		records = append(records, rec)
	}

	bites := tableme.TableMe(headers, records, noHeaders)
	util.PrintColorizedTable(bites, "asg", Config.Tableme.Colorize)
}

func setup() (*autoscaling.AutoScaling, error) {
	var client *autoscaling.AutoScaling = autoscaling.New(session.New(), aws.NewConfig())
	return client, nil
}
