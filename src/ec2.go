package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const region = "ap-northeast-1"

type Instance struct {
	InstanceId      string
	InstanceName    string
	State           string
	PublicIpAddress string
}

func getEC2Instances(profile string) ([]*Instance, error) {
	log.Println("getEC2Instances profile:", profile)
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(region)},
		Profile:           profile,
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		log.Println("Error creating session:", err)
		return nil, err
	}

	svc := ec2.New(sess)

	input := &ec2.DescribeInstancesInput{}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		log.Println("Error describing instances:", err)
		return nil, err
	}

	instances := []*Instance{}
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			instance := Instance{
				InstanceId:   *instance.InstanceId,
				InstanceName: getInstanceName(instance.Tags),
				State:        *instance.State.Name,
				PublicIpAddress: func() string {
					if instance.PublicIpAddress != nil {
						return *instance.PublicIpAddress
					} else {
						return ""
					}
				}(),
			}
			instances = append(instances, &instance)
		}
	}
	return instances, nil
}

func getInstanceName(tags []*ec2.Tag) string {
	for _, tag := range tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}
	return ""
}

func rebootInstance(instanceID string, profile string) {
	log.Println("rebootInstance")
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(region)},
		Profile:           profile,
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		log.Println("Error creating session:", err)
		return
	}

	svc := ec2.New(sess)
	input := &ec2.RebootInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	}

	_, err = svc.RebootInstances(input)
	if err != nil {
		log.Println("Error rebooting instance:", err)
		return
	}

	log.Println("Instance rebooted:", instanceID)
}

func stopInstance(instanceID string, profile string) {
	log.Println("stopInstance")
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(region)},
		Profile:           profile,
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		log.Println("Error creating session:", err)
		return
	}

	svc := ec2.New(sess)
	input := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	}

	_, err = svc.StopInstances(input)
	if err != nil {
		log.Println("Error stoping instance:", err)
		return
	}

	log.Println("Instance stopped:", instanceID)
}

func startInstance(instanceID string, profile string) {
	log.Println("startInstance")
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(region)},
		Profile:           profile,
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		log.Println("Error creating session:", err)
		return
	}

	svc := ec2.New(sess)
	input := &ec2.StartInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	}

	_, err = svc.StartInstances(input)
	if err != nil {
		log.Println("Error starting instance:", err)
		return
	}

	log.Println("Instance started:", instanceID)
}
