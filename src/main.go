package main

import (
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

const instanceTabName = "Instances  "
const settingsTabName = "Settings  "
const workspacesTabName = "Workspaces  "
const updateInterval = 20 * time.Second

func main() {
	settings, err := LoadSettings()
	if err != nil {
		log.Println("Error loading settings:", err)
		return
	}

	instances := []*Instance{}
	workspaces := []*Workspace{}
	instances = []*Instance{}
	workspaces = []*Workspace{}

	myApp := app.New()
	myWindow := myApp.NewWindow("EC2 Instances")

	ec2ListScreen := createEC2ListView(instances, "")
	workspaceListScreen := createWorkSpacesListView(workspaces, "")
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
		activeProfile := settings.GetActiveProfile()
		if item.Text == instanceTabName {
			if settings.GetActiveProfile() != "" {
				updateInstanceStatus(ec2ListScreen, activeProfile, activeFeature.EC2)
			}
		}
		if item.Text == workspacesTabName {
			if settings.GetActiveProfile() != "" {
				updateWorkspacesStatus(workspaceListScreen, activeProfile, activeFeature.WorkSpace)
			}
		}
	}

	myWindow.SetContent(ec2ListScreen)
	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(1200, 600))

	go updateScreen(settings, ec2ListScreen, workspaceListScreen)

	myWindow.ShowAndRun()
}

// 定期的にインスタンスの状態を更新する
func updateScreen(settings *Settings, ec2ListScreen *fyne.Container, workspaceListScreen *fyne.Container) {
	ticker := time.NewTicker(updateInterval)
	for range ticker.C {
		if settings.GetActiveProfile() != "" {
			activeFeature := settings.GetActiveFeature()
			updateInstanceStatus(ec2ListScreen, settings.GetActiveProfile(), activeFeature.EC2)
			updateWorkspacesStatus(workspaceListScreen, settings.GetActiveProfile(), activeFeature.WorkSpace)
		}
	}
}
