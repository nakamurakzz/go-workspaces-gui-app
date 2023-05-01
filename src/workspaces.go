package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/workspaces"
)

type Workspace struct {
	WorkspaceId string
	UserName    string
	State       string
}

func getWorkspaces(profile string) ([]*Workspace, error) {
	log.Println("getWorkspaces profile:", profile)
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(region)},
		Profile:           profile,
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		log.Println("Error creating session:", err)
		return nil, err
	}

	svc := workspaces.New(sess)

	input := &workspaces.DescribeWorkspacesInput{}

	result, err := svc.DescribeWorkspaces(input)
	if err != nil {
		log.Println("Error describing instances:", err)
		return nil, err
	}

	workspaceList := []*Workspace{}
	for _, workspaceInstance := range result.Workspaces {
		workspace := Workspace{
			WorkspaceId: *workspaceInstance.WorkspaceId,
			UserName:    *workspaceInstance.UserName,
			State:       *workspaceInstance.State,
		}
		workspaceList = append(workspaceList, &workspace)
	}
	return workspaceList, nil
}

func rebootWorkspce(workspaceID string, profile string) {
	log.Println("rebootWorkspce")
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(region)},
		Profile:           profile,
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		log.Println("Error creating session:", err)
		return
	}

	svc := workspaces.New(sess)
	input := &workspaces.RebootWorkspacesInput{
		RebootWorkspaceRequests: []*workspaces.RebootRequest{},
	}

	// _, err = svc.RebootWorkspaces(input)
	log.Println(svc, input)
	if err != nil {
		log.Println("Error rebooting instance:", err)
		return
	}

	log.Println("Instance rebooted:", workspaceID)
}
