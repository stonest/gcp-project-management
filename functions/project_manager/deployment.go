package deploymenthandler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"google.golang.org/api/deploymentmanager/v2"
)

//Parent details the parent container propeties for a project.
type Parent struct {
	Type string `json:"type"` // Type: Parent container type, valid values are organization or folder.
	ID   string `json:"id"`   // ID: The ID of the container.
}

//ProjectProperties provides information describing a project.
type ProjectProperties struct {
	Name      string `json:"name"`      // Name: The name of the project to deployed.
	ProjectID string `json:"projectId"` // ProjectID: The ID of the project to deployed.
	Parent    Parent `json:"parent"`    // Parent: Parent Container Properties.
}

//BillingProperties provides information describing a billing account.
type BillingProperties struct {
	Name             string `json:"name"`               // Name: Name of the project to assign billing to.
	BillingAccountID string `json:"billingAccountName"` // BillingAccountID: The name of the billing account as represented by its GUID. e.g. 014289-16D89B-9464F5.
}

//Resource builds out a GCP resource.
type Resource struct {
	Name       string      `json:"name"`               // Name: The name of the resource.
	Type       string      `json:"type"`               // Type: The type of GCP resource that is being deployed.
	Properties interface{} `json:"properties"`         // Properties: The properties of the resource to be deployed. This varies per resource.
	Metadata   Metadata    `json:"metadata,omitempty"` // Metadata: Optional metadata to include with the resource.
}

//Metadata provides additional metadata for a resource.
type Metadata struct {
	DependsOn     []string `json:"dependsOn,omitempty"`     // DependsOn: Specify another resource that this resource depends on.
	RuntimePolicy []string `json:"runtimePolicy,omitempty"` // RuntimePolicy: whatever google says it is.
}

//ProjectDeployment contains information for inserting, patching and deleting deployments.
type ProjectDeployment struct {
	Name           string `json:"name"`                     // Name: Name of the project requested for deployment.
	BillingAccount string `json:"billingAccount,omitempty"` // BillingAccount: The ID of the billing account to link the project to.
	ParentID       string `json:"parentId,omitempty"`       // ParentID: Parent container ID for the project.
	ParentType     string `json:"parentType,omitempty"`     // Type: Parent container type, valid values are organization or folder.
	Owner          string `json:"owner,omitempty"`          // Owner: Username that will own the project when complete.
}

//Resources is a container for n amount of Resource types.
type Resources struct {
	Resources []Resource `json:"resources"` // Resources: Container for Resources to be deployed via Deployment Manager.
}

//Insert will Insert a new GCP deployment of a new project.
func (projectDeployment *ProjectDeployment) Insert(ctx context.Context) (string, int, error) {
	resources := Resources{
		Resources: []Resource{
			Resource{
				Name: "project_" + projectDeployment.Name,
				Type: "cloudresourcemanager.v1.project",
				Properties: ProjectProperties{
					Name:      projectDeployment.Name,
					ProjectID: projectDeployment.Name,
					Parent: Parent{
						Type: projectDeployment.ParentType,
						ID:   projectDeployment.ParentID,
					},
				},
			},
			Resource{
				Name: "billing_" + projectDeployment.Name,
				Type: "deploymentmanager.v2.virtual.projectBillingInfo",
				Properties: BillingProperties{
					Name:             "projects/" + projectDeployment.Name,
					BillingAccountID: "billingAccounts/" + projectDeployment.BillingAccount,
				},
				Metadata: Metadata{
					DependsOn: []string{"project_" + projectDeployment.Name},
				},
			},
		},
	}
	deploymentConfig, _ := json.Marshal(resources)

	deployment := deploymentmanager.Deployment{
		Name: "deployment-" + projectDeployment.Name,
		Target: &deploymentmanager.TargetConfiguration{
			Config: &deploymentmanager.ConfigFile{
				Content: string(deploymentConfig),
			},
		},
	}
	resp, err := deploymentmanagerService.Deployments.Insert(projectID, &deployment).Context(ctx).Do()
	if err != nil {
		return "Error creating deployment", http.StatusInternalServerError, err
	}
	err = getDeploymentStatus(ctx, resp.Name)
	if err != nil {
		return "Error deploying project " + projectDeployment.Name, http.StatusInternalServerError, err
	}
	return "Successfully created " + projectDeployment.Name, http.StatusOK, nil
}

//Checks the operation deployment and returns the status of the deployment once the operation is complete.
func getDeploymentStatus(ctx context.Context, operation string) error {
	getResponse := deploymentmanagerService.Operations.Get(projectID, operation).Context(ctx)
	for {
		resp, err := getResponse.Do()
		if resp.Status == "DONE" {
			if resp.Error != nil {
				responseError, _ := resp.Error.MarshalJSON()
				return errors.New(string(responseError))
			}
			return nil
		}
		if err != nil {
			return err
		}
		log.Println("Waiting for deployment to complete...")
	}
}
