package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func createWorkSpacesListView(workspaces []*Workspace, profile string) *fyne.Container {
	description := widget.NewLabel("No profile selected")
	if profile != "" {
		description = widget.NewLabel("Workspaces in profile: " + profile)
	}
	content := container.NewVBox(description)
	rows := createWorkspacesList(content, workspaces, profile)
	for _, row := range rows {
		content.Add(row)
	}

	return content
}

func createWorkspacesList(content *fyne.Container, workspaces []*Workspace, profile string) []*fyne.Container {
	// Create header
	header := container.NewGridWithColumns(7,
		widget.NewLabel("Name"),
		widget.NewLabel("Instance ID"),
		widget.NewLabel("Status"),
		widget.NewLabel("Reboot"),
	)

	// Create rows
	rows := []*fyne.Container{header}
	for _, workspace := range workspaces {
		workspaceID := workspace.WorkspaceId
		userName := workspace.UserName
		status := workspace.State

		workspaceLabel := widget.NewLabel(workspaceID)
		statusLabel := widget.NewLabel(status)
		userNameLabel := widget.NewLabel(userName)
		rebootButton := widget.NewButton("Reboot", func() {
			rebootWorkspce(workspaceID, profile)
			updateWorkspacesStatus(content, profile)
		})

		if status == "stopped" {
			rebootButton.Disable()
		}
		if status == "pending" || status == "stopping" || status == "shutting-down" || status == "terminated" {
			rebootButton.Disable()
		}
		row := container.NewGridWithColumns(4, userNameLabel, workspaceLabel, statusLabel, rebootButton)
		rows = append(rows, row)
	}
	return rows
}

func updateWorkspacesStatus(content *fyne.Container, profile string) error {
	log.Println("Updating workspaces status")
	workspaces, err := getWorkspaces(profile)
	if err != nil {
		log.Println("Error:", err)
		return err
	}

	if profile != "" {
		// Create new rows based on the new instances
		newRows := createWorkspacesList(content, workspaces, profile)
		// Update content container
		content.Objects = nil
		for _, row := range newRows {
			content.Add(row)
		}
	}
	content.Refresh()

	return nil
}
