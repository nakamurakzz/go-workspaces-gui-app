package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func createEC2ListView(instances []*Instance, profile string) *fyne.Container {
	description := widget.NewLabel("No profile selected")
	if profile != "" {
		description = widget.NewLabel("EC2 Instances in profile: " + profile)
	}
	content := container.NewVBox(description)
	rows := createInstanceList(content, instances, profile)
	for _, row := range rows {
		content.Add(row)
	}

	return content
}

func createInstanceList(content *fyne.Container, instances []*Instance, profile string) []*fyne.Container {
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
		instanceID := instance.InstanceId
		status := instance.State
		ip := instance.PublicIpAddress
		name := instance.InstanceName

		nameLabel := widget.NewLabel(name)
		instanceLabel := widget.NewLabel(instanceID)
		ipLabel := widget.NewLabel(ip)
		statusLabel := widget.NewLabel(status)
		rebootButton := widget.NewButton("Reboot", func() {
			rebootInstance(instanceID, profile)
			updateInstanceStatus(content, profile, true)
		})
		startButton := widget.NewButton("Start", func() {
			startInstance(instanceID, profile)
			updateInstanceStatus(content, profile, true)
		})
		stopButton := widget.NewButton("Stop", func() {
			stopInstance(instanceID, profile)
			updateInstanceStatus(content, profile, true)
		})

		if status == "running" {
			startButton.Disable()
		} else if status == "stopped" {
			stopButton.Disable()
			rebootButton.Disable()
		}
		if status == "pending" || status == "stopping" || status == "shutting-down" || status == "terminated" {
			rebootButton.Disable()
			startButton.Disable()
			stopButton.Disable()
		}
		row := container.NewGridWithColumns(7, nameLabel, instanceLabel, ipLabel, statusLabel, rebootButton, startButton, stopButton)
		rows = append(rows, row)
	}
	return rows
}

func updateInstanceStatus(content *fyne.Container, profile string, isActive bool) error {
	if !isActive {
		content.Refresh()
		return nil
	}
	log.Println("updateInstanceStatus")
	instances, err := getEC2Instances(profile)
	if err != nil {
		log.Println("Error:", err)
		return err
	}

	if profile != "" {
		// Create new rows based on the new instances
		newRows := createInstanceList(content, instances, profile)
		// Update content container
		content.Objects = nil
		for _, row := range newRows {
			content.Add(row)
		}
	}
	content.Refresh()

	return nil
}
