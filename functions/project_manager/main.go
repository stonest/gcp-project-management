package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/deploymentmanager/v2"
)

var projectID = os.Getenv("GCP_PROJECT")
var deploymentmanagerService *deploymentmanager.Service
var cloudresourcemanagerService *cloudresourcemanager.Service

// Set context and create a new deployment manager service that will persist between runs.
func init() {
	var err error
	deploymentmanagerService, err = deploymentmanager.NewService(context.Background())
	if err != nil {
		log.Fatalf("deploymentManager.NewService: %v", err)
	}
	cloudresourcemanagerService, err = cloudresourcemanager.NewService(context.Background())
	if err != nil {
		log.Fatalf("cloudresourcemanagerService.NewService: %v", err)
	}
}

func main() {
	// Basic handler for now
	http.Handle("/", deploymentHandler(deployment))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
