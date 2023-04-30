package main

import (
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const instanceTabName = "Instances  "
const settingsTabName = "Settings  "
const workspacesTabName = "Workspaces  "

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

func main() {
	settings, err := LoadSettings()
	if err != nil {
		log.Println("Error loading settings:", err)
		return
	}

	instances := []*Instance{}
	workspaces := []*Workspace{}
	if settings.GetActiveProfile() == "" {
		instances = []*Instance{}
		workspaces = []*Workspace{}
	} else {
		instances, err = getEC2Instances(settings.GetActiveProfile())
		if err != nil {
			log.Println("Error:", err)
			return
		}
		workspaces, err = getWorkspaces(settings.GetActiveProfile())
		if err != nil {
			log.Println("Error:", err)
			return
		}
	}

	myApp := app.New()
	myWindow := myApp.NewWindow("EC2 Instances")

	ec2ListScreen := createEC2ListView(instances, settings.GetActiveProfile())
	workspaceListScreen := createWorkSpacesListView(workspaces, settings.GetActiveProfile())
	settingsScreen := createSettingsScreen(myApp, settings)
	tabs := container.NewAppTabs(
		container.NewTabItem(settingsTabName, settingsScreen),
		container.NewTabItem(instanceTabName, ec2ListScreen),
		container.NewTabItem(workspacesTabName, workspaceListScreen),
	)
	tabs.SetTabLocation(container.TabLocationTop)
	tabs.SelectIndex(0)

	tabs.OnSelected = func(item *container.TabItem) {
		activeFeature := settings.GetActiveFeature()
		if item.Text == instanceTabName {
			if settings.GetActiveProfile() != "" {
				updateInstanceStatus(ec2ListScreen, settings.GetActiveProfile(), activeFeature.EC2)
			}
		}
		if item.Text == workspacesTabName {
			if settings.GetActiveProfile() != "" {
				updateWorkspacesStatus(workspaceListScreen, settings.GetActiveProfile(), activeFeature.WorkSpace)
			}
		}
	}

	myWindow.SetContent(ec2ListScreen)
	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(1200, 600))

	ticker := time.NewTicker(20 * time.Second)
	go func(settings *Settings) {
		for range ticker.C {
			if settings.GetActiveProfile() != "" {
				activeFeature := settings.GetActiveFeature()
				updateInstanceStatus(ec2ListScreen, settings.GetActiveProfile(), activeFeature.EC2)
				updateWorkspacesStatus(workspaceListScreen, settings.GetActiveProfile(), activeFeature.WorkSpace)
			}
		}
	}(settings)

	myWindow.ShowAndRun()
}
