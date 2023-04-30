package main

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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

func createStatusLabel(status string) *widget.Label {
	// var textColor color.Color
	// switch status {
	// case "running":
	// 	textColor = color.RGBA{0, 255, 0, 255}
	// case "stopped":
	// 	textColor = color.RGBA{255, 0, 0, 255}
	// default:
	// 	textColor = color.RGBA{128, 128, 128, 255}
	// }
	return widget.NewLabelWithStyle(status, fyne.TextAlignLeading, fyne.TextStyle{})
}

func updateInstanceStatus(rows []*fyne.Container) {
	instances, err := getEC2Instances()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for i, instance := range instances {
		status := *instance.State.Name
		row := rows[i]

		newStatusLabel := createStatusLabel(status)
		row.Objects[1] = newStatusLabel
		row.Refresh()

		if instance.PrivateIpAddress != nil {
			row.Objects[3].(*widget.Label).SetText(*instance.PrivateIpAddress)
		}
		if name := getInstanceName(instance.Tags); name != "" {
			row.Objects[4].(*widget.Label).SetText(name)
		}
	}
}

func getInstanceName(tags []*ec2.Tag) string {
	for _, tag := range tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}
	return ""
}

func main() {
	instances, err := getEC2Instances()
	if err != nil {
		log.Println("Error:", err)
		return
	}

	myApp := app.New()
	myWindow := myApp.NewWindow("EC2 Instances")

	// Create header
	header := container.NewGridWithColumns(7,
		widget.NewLabel("Name"),
		widget.NewLabel("Instance ID"),
		widget.NewLabel("IP Address"),
		widget.NewLabel("Status"),
		widget.NewLabel("Reboot"),
		widget.NewLabel("Start"),
		widget.NewLabel("Stop"),
	)

	// Create rows
	rows := []*fyne.Container{header}
	for _, instance := range instances {
		instanceID := *instance.InstanceId
		status := *instance.State.Name
		ip := ""
		if instance.PrivateIpAddress != nil {
			ip = *instance.PrivateIpAddress
		}
		name := getInstanceName(instance.Tags)

		nameLabel := widget.NewLabel(name)
		instanceLabel := widget.NewLabel(instanceID)
		ipLabel := widget.NewLabel(ip)
		statusLabel := widget.NewLabel(status)
		rebootButton := widget.NewButton("Reboot", func() {
			rebootInstance(instanceID)
			updateInstanceStatus(rows)
		})
		startButton := widget.NewButton("Start", func() {
			startInstance(instanceID)
			updateInstanceStatus(rows)
		})
		stopButton := widget.NewButton("Stop", func() {
			stopInstance(instanceID)
			updateInstanceStatus(rows)
		})

		row := container.NewGridWithColumns(7, nameLabel, instanceLabel, ipLabel, statusLabel, rebootButton, startButton, stopButton)
		rows = append(rows, row)
	}

	// Create content container
	content := container.NewVBox()
	for _, row := range rows {
		content.Add(row)
	}

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(800, 600))

	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			updateInstanceStatus(rows)
		}
	}()

	myWindow.ShowAndRun()
}
