package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func createInstanceList(content *fyne.Container, instances []*ec2.Instance, profile string) []*fyne.Container {
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
			rebootInstance(instanceID, profile)
			updateInstanceStatus(content, profile)
		})
		startButton := widget.NewButton("Start", func() {
			startInstance(instanceID, profile)
			updateInstanceStatus(content, profile)
		})
		stopButton := widget.NewButton("Stop", func() {
			stopInstance(instanceID, profile)
			updateInstanceStatus(content, profile)
		})

		row := container.NewGridWithColumns(7, nameLabel, instanceLabel, ipLabel, statusLabel, rebootButton, startButton, stopButton)
		rows = append(rows, row)
	}
	return rows
}
