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
	var deploymentErrorResponse *APIError

	newDeployment := ProjectDeployment{}
	data, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(data, &newDeployment)

	switch method := r.Method; method {
	case "POST":
		deploymentErrorResponse = newDeployment.Insert(r.Context())
	case "DELETE":
		deploymentErrorResponse = newDeployment.Delete(r.Context())
	}
	if deploymentErrorResponse != nil {
		log.Println(deploymentErrorResponse.Error.Error())
		http.Error(w, deploymentErrorResponse.Message, deploymentErrorResponse.Code)
		return
	}
	fmt.Fprintf(w, newDeployment.Name+" "+r.Method+" Successful")
}
