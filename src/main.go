package main

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/aws/aws-sdk-go/service/ec2"
)

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

	header := container.NewHBox(
		widget.NewLabel("Instance ID"),
		widget.NewLabel("Status"),
		widget.NewLabel("Reboot"),
		widget.NewLabel("IP Address"),
		widget.NewLabel("Name"),
	)
	rows := []*fyne.Container{header}

	for _, instance := range instances {
		instanceID := *instance.InstanceId
		status := *instance.State.Name
		ip := ""
		if instance.PrivateIpAddress != nil {
			ip = *instance.PrivateIpAddress
		}
		name := getInstanceName(instance.Tags)

		instanceLabel := widget.NewLabel(instanceID)
		statusLabel := widget.NewLabel(status)
		ipLabel := widget.NewLabel(ip)
		nameLabel := widget.NewLabel(name)
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

		row := container.NewHBox(nameLabel, instanceLabel, ipLabel, statusLabel, rebootButton, startButton, stopButton)
		rows = append(rows, row)
	}

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
