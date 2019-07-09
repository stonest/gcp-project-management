package deploymenthandler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"google.golang.org/api/deploymentmanager/v2"
)

var projectID = "org-dev"
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
	log.Println(r.Method)
	newDeployment := ProjectDeployment{}
	data, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(data, &newDeployment)
	response, status, err := newDeployment.Insert(r)
	if err != nil {
		http.Error(w, response, status)
		log.Print(err)
		return
	}
	fmt.Fprintf(w, response)
	return
}
