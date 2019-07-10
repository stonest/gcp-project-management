package deploymenthandler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/deploymentmanager/v2"
)

var projectID = os.Getenv("GCP_PROJECT")
var deploymentmanagerService *deploymentmanager.Service

// Set context and create a new deployment manager service that will perist between runs.
func init() {
	var err error
	// use context.Background() here to persist context between invocations.
	deploymentmanagerService, err = deploymentmanager.NewService(context.Background())
	if err != nil {
		log.Fatalf("deploymentManager.NewService: %v", err)
	}
}

// ManageDeployment Creates, updates and deletes project deployments.
func ManageDeployment(w http.ResponseWriter, r *http.Request) {
	var response string
	var status int
	var err error

	newDeployment := ProjectDeployment{}
	data, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(data, &newDeployment)

	switch method := r.Method; method {
	case "POST":
		response, status, err = newDeployment.Insert(r.Context())
	}
	if err != nil {
		http.Error(w, response, status)
		log.Print(err)
		return
	}
	fmt.Fprintf(w, response)
}
