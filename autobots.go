package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/codegangsta/cli"
)

func Assert(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		os.Exit(1)
	}
}

func GetAutoScalingGroup(instance, region string) *string {
	config := aws.NewConfig().WithRegion(region)
	svcEc := ec2.New(config)

	resp, err := svcEc.DescribeInstances(
		&ec2.DescribeInstancesInput{InstanceIds: []*string{&instance}},
	)

	Assert(err)

	for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			for _, tag := range inst.Tags {
				if *tag.Key == "aws:autoscaling:groupName" {
					return tag.Value
				}
			}
		}
	}
	return nil
}

func RollOut(c *configRollOut) {
	svcMeta := ec2metadata.New(&ec2metadata.Config{})

	region, err := svcMeta.Region()
	Assert(err)

	config := aws.NewConfig().WithRegion(region)
	svcAsg := autoscaling.New(config)
	svcEc := ec2.New(config)

	instanceId, err := svcMeta.GetMetadata("instance-id")
	instanceASG := GetAutoScalingGroup(instanceId, region)

	if instanceASG == nil && len(c.groups) == 0 {
		fmt.Fprintln(os.Stderr, "missing autoscaling groups")
		os.Exit(1)
	}

	autoScalingGroups := make([]*string, 0)

	if instanceASG != nil {
		autoScalingGroups = append(autoScalingGroups, instanceASG)
	}

	for _, item := range c.groups {
		autoScalingGroups = append(autoScalingGroups, aws.String(item))
	}

	asg, err := svcAsg.DescribeAutoScalingGroups(
		&autoscaling.DescribeAutoScalingGroupsInput{
			AutoScalingGroupNames: autoScalingGroups,
		},
	)

	Assert(err)

	instances := make([]*string, 0)
	for idx, _ := range asg.AutoScalingGroups {
		for _, inst := range asg.AutoScalingGroups[idx].Instances {
			instances = append(instances, inst.InstanceId)
		}
	}

	if len(instances) == 0 {
		os.Exit(1)
	}

	resp, err := svcEc.DescribeInstances(
		&ec2.DescribeInstancesInput{InstanceIds: instances},
	)

	Assert(err)

	values := make(results, 0)
	for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			value := ""
			switch c.filter {
			case "public-dns":
				value = *inst.PublicDnsName
			case "private-dns":
				value = *inst.PrivateDnsName
			case "hostname":
				value = fmt.Sprintf("ip-%s", strings.Replace(*inst.PrivateIpAddress, ".", "-", -1))
			default:
				value = *inst.PrivateIpAddress
			}
			if value != "" {
				values = append(values, result{instance: value, launchTime: *inst.LaunchTime})
			}
		}
	}

	if len(values) > 0 {
		sort.Sort(values)
		fmt.Println(strings.Join(values.Values(), " "))
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "autobots"
	app.Version = "v1.0.0"
	app.Usage = "autobots assemble!"
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "with-asg",
			Usage: "auto scaling groups to list",
		},
		cli.StringFlag{
			Name:  "output",
			Usage: "addresses format: private-ip|private-dns|public-dns|hostname",
			Value: "private-ip",
		},
	}

	app.Action = func(c *cli.Context) {
		groups := c.StringSlice("with-asg")
		filter := c.String("output")

		switch filter {
		case "private-ip":
			//
		case "public-dns":
			//
		case "private-dns":
			//
		case "hostname":
			//
		default:
			fmt.Println("invalid output")
			os.Exit(1)
		}

		// autobots roll out!
		RollOut(&configRollOut{
			groups: groups,
			filter: filter,
		},
		)
	}

	app.Run(os.Args)
}
