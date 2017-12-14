package ec2

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	config "github.com/alexebird/aws-find/config"
	util "github.com/alexebird/aws-find/util"
	"github.com/alexebird/tableme/tableme"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/kballard/go-shellquote"
	//"github.com/davecgh/go-spew/spew"
)

var Config config.AwsFindConfig

func filterInstances(instances []*ec2.Instance, test func(*ec2.Instance) bool) (ret []*ec2.Instance) {
	for _, i := range instances {
		if test(i) {
			ret = append(ret, i)
		}
	}
	return
}

func describeInstances(client *ec2.EC2, all bool, nameFilter string) []*ec2.Instance {
	var filters []*ec2.Filter

	if all == true {
		filters = nil
	} else {
		filters = []*ec2.Filter{
			{
				Name:   aws.String("tag:env"),
				Values: []*string{aws.String("dev-0")},
			},
		}
	}

	params := &ec2.DescribeInstancesInput{Filters: filters}
	instances := make([]*ec2.Instance, 0)

	err := client.DescribeInstancesPages(params,
		func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
			for _, res := range page.Reservations {
				instances = append(instances, res.Instances...)
			}
			return true
		})

	if err != nil {
		panic(err)
	}

	if len(nameFilter) > 0 {
		filterTest := func(i *ec2.Instance) bool {
			nameTagValue := findTagByKey(i, "Name")
			if nameTagValue != nil {
				return strings.Contains(*nameTagValue, nameFilter)
			} else {
				return false
			}
		}
		instances = filterInstances(instances, filterTest)
	}

	return instances
}

func findTagByKey(instance *ec2.Instance, key string) *string {
	for _, tag := range instance.Tags {
		if *tag.Key == key {
			return tag.Value
		}
	}

	return nil
}

func davinciShortFormTable(instances []*ec2.Instance) {
	headers := []string{
		"PRIVATE_IP", "NAME", "COLOR", "ROLE", "STATE", "TYPE", "IMAGE", "KEY",
	}

	records := make([][]string, 0)

	for _, inst := range instances {
		rec := []string{
			*inst.PrivateIpAddress,
			*findTagByKey(inst, "Name"),
			*findTagByKey(inst, "color"),
			*findTagByKey(inst, "role"),
			*inst.State.Name,
			*inst.InstanceType,
			*inst.ImageId,
			*inst.KeyName,
		}
		records = append(records, rec)
	}

	bites := tableme.TableMe(headers, records)
	util.PrintColorizedTable(bites, "ec2", Config.Tableme.Colorize)
}

func davinciLongFormTable(instances []*ec2.Instance) {
	headers := []string{
		"PUBLIC_IP", "PRIVATE_IP", "NAME", "COLOR", "ROLE", "ENV", "STATE", "TYPE", "IMAGE", "LAUNCHED", "KEY", "ID", //"SUBNET", "CIDR",
	}

	records := make([][]string, 0)

	for _, inst := range instances {
		time := inst.LaunchTime.Format(time.RFC3339)

		rec := []string{
			tableme.WithEmptyStringDefault(inst.PublicIpAddress),
			tableme.WithEmptyStringDefault(inst.PrivateIpAddress),
			tableme.WithEmptyStringDefault(findTagByKey(inst, "Name")),
			tableme.WithEmptyStringDefault(findTagByKey(inst, "color")),
			tableme.WithEmptyStringDefault(findTagByKey(inst, "role")),
			tableme.WithEmptyStringDefault(findTagByKey(inst, "env")),
			tableme.WithEmptyStringDefault(inst.State.Name),
			tableme.WithEmptyStringDefault(inst.InstanceType),
			tableme.WithEmptyStringDefault(inst.ImageId),
			tableme.WithEmptyStringDefault(&time),
			tableme.WithEmptyStringDefault(inst.KeyName),
			tableme.WithEmptyStringDefault(inst.InstanceId),
		}
		records = append(records, rec)
	}

	bites := tableme.TableMe(headers, records)
	util.PrintColorizedTable(bites, "ec2", Config.Tableme.Colorize)
}

func shortFormTable(instances []*ec2.Instance) {
	headers := []string{
		"PRIVATE_IP", "NAME", "STATE", "TYPE", "IMAGE", "KEY",
	}

	records := make([][]string, 0)

	for _, inst := range instances {
		rec := []string{
			*inst.PrivateIpAddress,
			*findTagByKey(inst, "Name"),
			*inst.State.Name,
			*inst.InstanceType,
			*inst.ImageId,
			*inst.KeyName,
		}
		records = append(records, rec)
	}

	bites := tableme.TableMe(headers, records)
	util.PrintColorizedTable(bites, "ec2", Config.Tableme.Colorize)
}

func longFormTable(instances []*ec2.Instance) {
	headers := []string{
		"PUBLIC_IP", "PRIVATE_IP", "NAME", "STATE", "TYPE", "IMAGE", "LAUNCHED", "KEY", "ID", //"SUBNET", "CIDR",
	}

	records := make([][]string, 0)

	for _, inst := range instances {
		time := inst.LaunchTime.Format(time.RFC3339)

		rec := []string{
			tableme.WithEmptyStringDefault(inst.PublicIpAddress),
			tableme.WithEmptyStringDefault(inst.PrivateIpAddress),
			tableme.WithEmptyStringDefault(findTagByKey(inst, "Name")),
			tableme.WithEmptyStringDefault(inst.State.Name),
			tableme.WithEmptyStringDefault(inst.InstanceType),
			tableme.WithEmptyStringDefault(inst.ImageId),
			tableme.WithEmptyStringDefault(&time),
			tableme.WithEmptyStringDefault(inst.KeyName),
			tableme.WithEmptyStringDefault(inst.InstanceId),
		}
		records = append(records, rec)
	}

	bites := tableme.TableMe(headers, records)
	util.PrintColorizedTable(bites, "ec2", Config.Tableme.Colorize)
}

func sshCommand(ip string) []string {
	return []string{"ssh", ip}
}

func execInteractive(cmdArgs []string) {
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:len(cmdArgs)]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func chooseInstanceForConnect(instances []*ec2.Instance) *ec2.Instance {
	return instances[0]
}

func connect(inst *ec2.Instance) {
	name := tableme.WithEmptyStringDefault(findTagByKey(inst, "Name"))
	ip := *inst.PrivateIpAddress

	cmd := sshCommand(ip)

	fmt.Println()
	fmt.Printf("==> connecting to %s(%s)\n", name, ip)
	fmt.Printf("==> via command: %s\n", shellquote.Join(cmd...))
	fmt.Println()

	execInteractive(cmd)
}

func setup() (*ec2.EC2, error) {
	var client *ec2.EC2 = ec2.New(session.New())
	return client, nil
}
