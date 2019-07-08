package main

import (
	"encoding/json"
	"fmt"
	"log"

	"google.golang.org/api/deploymentmanager/v2"
)

//Parent details the parent contianer propeties for a project
type Parent struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

//ProjectProperties provides information describing a project
type ProjectProperties struct {
	Name      string `json:"name"`
	ProjectID string `json:"projectId"`
	Parent    Parent `json:"parent"`
}

//BillingProperties provides information describing a billing account
type BillingProperties struct {
	Name               string `json:"name"`
	BillingAccountName string `json:"billingAccountName"`
}

//Resource builds out a GCP resource.
type Resource struct {
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Properties interface{} `json:"properties"`
	Metadata   Metadata    `json:"metadata,omitempty"`
}

//Metadata provides additional metadata for a resource
type Metadata struct {
	DependsOn     []string `json:"dependsOn"`
	RuntimePolicy []string `json:"runtimePolicy,omitempty"`
}

//ProjectDeployment contains information for inserting, patching and deleting deployments.
type ProjectDeployment struct {
	Name           string `json:"name"`
	BillingAccount string `json:"billingAccount"`
	ParentID       string `json:"parentId"`
	ParentType     string `json:"parentType"`
	Owner          string `json:"owner"`
}

//Resources is a container for n amount of Resource types.
type Resources struct {
	Resources []Resource `json:"resources"`
}

//Insert will Insert a new GCP deployment of a new project.
func (projectDeployment *ProjectDeployment) Insert(service *deploymentmanager.Service) {
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
					Name:               "projects/" + projectDeployment.Name,
					BillingAccountName: "billingAccounts/" + projectDeployment.BillingAccount,
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
	resp, err := service.Deployments.Insert("org-dev", &deployment).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", resp)
}
