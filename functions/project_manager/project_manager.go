package deploymenthandler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"google.golang.org/api/deploymentmanager/v2"
)

var projectID = "org-dev"
var deploymentmanagerService *deploymentmanager.Service
var ctx context.Context

// Set context and create a new deployment manager service that will perist between runs.
func init() {
	var err error
	ctx = context.Background()
	deploymentmanagerService, err = deploymentmanager.NewService(ctx)
	if err != nil {
		log.Fatalf("deploymentManager.NewService: %v", err)
	}
}

// ManageDeployment Creates, updates and deletes project deployments.
func ManageDeployment(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method)
	newDeployment := ProjectDeployment{}
	data, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(data, &newDeployment)
	newDeployment.Insert()
}
