package deploymenthandler

import (
	"context"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/deploymentmanager/v2"
)

var projectID = os.Getenv("GCP_PROJECT")
var deploymentmanagerService *deploymentmanager.Service

// Set context and create a new deployment manager service that will persist between runs.
func init() {
	var err error

	apiKey, set := os.LookupEnv("API_KEY")
	if set == false {
		deploymentmanagerService, err = deploymentmanager.NewService(context.Background())
	} else {
		deploymentmanagerService, err = deploymentmanager.NewService(context.Background(), option.WithAPIKey(apiKey))
	}
	if err != nil {
		log.Fatalf("deploymentManager.NewService: %v", err)
	}
}

func main() {
	// Basic handler for now
	http.Handle("/", deploymentHandler(deployment))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
