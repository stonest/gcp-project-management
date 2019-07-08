package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/api/deploymentmanager/v2"
)

var projectID = os.Getenv("GCP_PROJECT")
var deploymentmanagerService *deploymentmanager.Service
var ctx context.Context

func init() {
	var err error
	ctx = context.Background()
	deploymentmanagerService, err = deploymentmanager.NewService(ctx)
	if err != nil {
		log.Fatalf("deploymentManager.NewService: %v", err)
	}
}

func main() {

	newDeployment := ProjectDeployment{
		BillingAccount: "014289-16D89B-9466F5",
		Name:           "proiecjts-sstone",
		ParentID:       "283639071922",
		ParentType:     "organization",
	}

	newDeployment.Insert(deploymentmanagerService)
}
