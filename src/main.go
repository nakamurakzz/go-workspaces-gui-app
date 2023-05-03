package main

import (
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

const (
	instanceTabName   = "EC2 Instances  "
	settingsTabName   = "Settings  "
	workspacesTabName = "Workspaces  "
	updateInterval    = 20 * time.Second
)

func main() {
	settings, err := LoadSettings()
	if err != nil {
		log.Println("Error loading settings:", err)
		return
	}

	instances := []*Instance{}
	workspaces := []*Workspace{}

	myApp := app.New()

	ec2List := createEC2ListView(instances, "")
	workspaceList := createWorkSpacesListView(workspaces, "")
	settingsScreen := createSettingsScreen(myApp, settings)

	tabs := container.NewAppTabs(
		container.NewTabItem(settingsTabName, settingsScreen),
		container.NewTabItem(instanceTabName, ec2List),
		container.NewTabItem(workspacesTabName, workspaceList),
	)
	tabs.SetTabLocation(container.TabLocationTop)
	tabs.SelectIndex(0)

	tabs.OnSelected = func(item *container.TabItem) {
		activeFeature := settings.GetActiveFeature()
		activeProfile := settings.GetActiveProfile()
		if item.Text == instanceTabName && activeProfile != "" {
			go updateInstanceStatus(ec2List, activeProfile, activeFeature.EC2)
		}
		if item.Text == workspacesTabName && activeProfile != "" {
			go updateWorkspacesStatus(workspaceList, activeProfile, activeFeature.WorkSpace)
		}
	}

	window := myApp.NewWindow("EC2 Instances")
	window.Resize(fyne.NewSize(1200, 600))
	window.SetContent(ec2List)
	window.SetContent(tabs)

	go updateScreen(settings, ec2List, workspaceList)

	window.ShowAndRun()
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
