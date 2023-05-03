package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

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

	ec2FeatureCheckbox := widget.NewCheck("Enable EC2 List Screen", func(value bool) {
		settings.changeEC2ListScreenVisible()
	})

	worksoaceFeatureCheckbox := widget.NewCheck("Enable Workspaces List Screen", func(value bool) {
		settings.changeWorkSpaceListScreenVisible()
	})

	description := widget.NewLabel("Enter AWS profile names.")
	featureSelectDescription := widget.NewLabel("Select Features to enable. When you enable a feature, it will start to fetch data from AWS.")
	settingsContainer := container.NewVBox(description, settingsForm, saveButton, profile, featureSelectDescription, ec2FeatureCheckbox, worksoaceFeatureCheckbox)
	return settingsContainer
}
