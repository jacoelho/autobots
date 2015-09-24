package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/codegangsta/cli"
)

type configRollOut struct {
	groups []string
	filter string
}

func Assert(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		os.Exit(1)
	}
}

func RollOut(c *configRollOut) {
	svcMeta := ec2metadata.New(&ec2metadata.Config{})

	region, err := svcMeta.Region()
	Assert(err)

	config := aws.NewConfig().WithRegion(region)
	svcAsg := autoscaling.New(config)
	svcEc := ec2.New(config)

	autoScalingGroups := make([]*string, len(c.groups))
	for idx, item := range c.groups {
		autoScalingGroups[idx] = aws.String(item)
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

	result := make([]string, 0)
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
			result = append(result, value)
		}
	}

	fmt.Println(strings.Join(result, " "))
}

func main() {
	app := cli.NewApp()
	app.Name = "autobots"
	app.Version = "v1.0.0"
	app.Usage = "autobots assemble!"
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "auto-scaling-groups",
			Usage: "auto scaling groups to list",
		},
		cli.StringFlag{
			Name:  "output",
			Usage: "addresses format: private-ip|private-dns|public-dns|hostname",
			Value: "private-ip",
		},
	}

	app.Action = func(c *cli.Context) {
		groups := c.StringSlice("auto-scaling-groups")
		filter := c.String("output")

		if !c.GlobalIsSet("auto-scaling-groups") {
			fmt.Println("need at least 1 auto-scaling-group")
			os.Exit(1)
		}

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
