package ec2

import (
	//"fmt"
	"os"

	"github.com/alexebird/tableme/tableme"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/davecgh/go-spew/spew"
)

//type Ec2Instance struct {
//}

func describeInstances(client *ec2.EC2, all bool) []*ec2.Instance {
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
				spew.Dump(res)
				instances = append(instances, res.Instances...)
			}
			return true
		})

	if err != nil {
		panic(err)
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
		spew.Dump(inst)
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

	err := tableme.TableMe(headers, records)
	if err != nil {
		os.Exit(1)
	}
}

func davinciLongFormTable(instances []*ec2.Instance) {
	headers := []string{
		"PUBLIC_IP", "PRIVATE_IP", "NAME", "COLOR", "ROLE", "ENV", "STATE", "TYPE", "IMAGE", "LAUNCHED", //"KEY", "SUBNET", "CIDR", "ID",
	}

	records := make([][]*string, 0)

	for _, inst := range instances {
		spew.Dump(inst)
		rec := 

		rec := []*string{
			*inst.PublicIpAddress,
			*inst.PrivateIpAddress,
			*findTagByKey(inst, "Name"),
			*findTagByKey(inst, "color"),
			*findTagByKey(inst, "role"),
			*findTagByKey(inst, "env"),
			*inst.State.Name,
			*inst.InstanceType,
			*inst.ImageId,
			*inst.KeyName,
		}
		records = append(records, rec)
	}

	err := tableme.TableMePtr(headers, records)
	if err != nil {
		os.Exit(1)
	}
}

func shortFormTable(instances []*ec2.Instance) {
	headers := []string{
		"PRIVATE_IP", "NAME", "STATE", "TYPE", "IMAGE", "KEY",
	}

	records := make([][]string, 0)

	for _, inst := range instances {
		spew.Dump(inst)
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

	err := tableme.TableMe(headers, records)
	if err != nil {
		os.Exit(1)
	}
}

func setup() (*ec2.EC2, error) {
	var client *ec2.EC2 = ec2.New(session.New())
	return client, nil
}
