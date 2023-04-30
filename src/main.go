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

const instanceTabName = "Instances  "
const settingsTabName = "Settings  "

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

func updateInstanceStatus(content *fyne.Container, profile string) error {
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

func getInstanceName(tags []*ec2.Tag) string {
	for _, tag := range tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}
	return ""
}

func main() {
	settings, err := LoadSettings()
	if err != nil {
		log.Println("Error loading settings:", err)
		return
	}

	instances := []*ec2.Instance{}

	if settings.GetActiveProfile() == "" {
		instances = []*ec2.Instance{}
	} else {
		instances, err = getEC2Instances(settings.GetActiveProfile())
		if err != nil {
			log.Println("Error:", err)
			return
		}
	}

	myApp := app.New()
	myWindow := myApp.NewWindow("EC2 Instances")

	// Create content container
	description := widget.NewLabel("EC2 Instances in profile: " + settings.GetActiveProfile())
	content := container.NewVBox(description)
	rows := createInstanceList(content, instances, settings.GetActiveProfile())
	for _, row := range rows {
		content.Add(row)
	}

	settingsScreen := createSettingsScreen(myApp, settings)
	tabs := container.NewAppTabs(
		container.NewTabItem(instanceTabName, content),
		container.NewTabItem(settingsTabName, settingsScreen),
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	tabs.OnSelected = func(item *container.TabItem) {
		if item.Text == instanceTabName {
			updateInstanceStatus(content, settings.GetActiveProfile())
		}
	}

	myWindow.SetContent(content)
	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(800, 600))

	ticker := time.NewTicker(20 * time.Second)
	go func() {
		for range ticker.C {
			updateInstanceStatus(content, settings.GetActiveProfile())
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

	profileRadio := widget.NewRadioGroup(settings.Profiles, func(value string) {
		settings.SetActiveProfile(value)
	})
	settingsForm.Append("Active Profile", profileRadio)

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
	settingsContent := container.NewVBox(description, settingsForm, saveButton)
	return settingsContent
}
