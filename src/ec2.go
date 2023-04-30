package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const region = "ap-northeast-1"

func getEC2Instances() ([]*ec2.Instance, error) {
	log.Println("getEC2Instances")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return nil, err
	}

	svc := ec2.New(sess)

	input := &ec2.DescribeInstancesInput{}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		return nil, err
	}

	instances := []*ec2.Instance{}
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			instances = append(instances, instance)
		}
	}

	return instances, nil
}

func rebootInstance(instanceID string) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
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

func stopInstance(instanceID string) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
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

func startInstance(instanceID string) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
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
