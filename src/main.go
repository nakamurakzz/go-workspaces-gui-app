package main

import (
	"fmt"
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
		if item.Text == instanceTabName {
			if settings.GetActiveProfile() != "" {
				updateInstanceStatus(ec2ListScreen, settings.GetActiveProfile())
			}
		}
		if item.Text == workspacesTabName {
			if settings.GetActiveProfile() != "" {
				updateWorkspacesStatus(workspaceListScreen, settings.GetActiveProfile())
			}
		}
	}

	myWindow.SetContent(ec2ListScreen)
	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(1200, 600))

	ticker := time.NewTicker(20 * time.Second)
	go func() {
		for range ticker.C {
			if settings.GetActiveProfile() != "" {
				updateInstanceStatus(ec2ListScreen, settings.GetActiveProfile())
				updateWorkspacesStatus(workspaceListScreen, settings.GetActiveProfile())
			}
		}
	}()

	myWindow.ShowAndRun()
}

func createSettingsScreen(app fyne.App, settings *Settings) fyne.CanvasObject {
	profileEntries := []*widget.Entry{}
	settingsForm := &widget.Form{}

	for i := 0; i < 5; i++ {
		entry := widget.NewEntry()

		if i < len(settings.Profiles) {
			entry.SetText(settings.Profiles[i])
		}

		profileEntries = append(profileEntries, entry)
		settingsForm.Append(fmt.Sprintf("Profile %d", i+1), entry)
	}

	profile := &widget.Form{}
	profileRadio := widget.NewRadioGroup(settings.Profiles, func(value string) {
		settings.SetActiveProfile(value)
	})
	profile.Append("Active Profile", profileRadio)

	saveButton := widget.NewButton("Save", func() {
		newProfiles := []string{}

		for _, entry := range profileEntries {
			profileName := entry.Text

			if profileName != "" {
				newProfiles = append(newProfiles, profileName)
			}
		}

		settings.Profiles = newProfiles
		profileRadio.Options = newProfiles
		profileRadio.Refresh()

		if err := settings.Save(); err != nil {
			log.Println("Error saving settings:", err)
		}
	})

	description := widget.NewLabel("Enter AWS profile names.")
	settingsContent := container.NewVBox(description, settingsForm, saveButton, profile)
	return settingsContent
}
