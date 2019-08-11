package deploymenthandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type deploymentHandler func(http.ResponseWriter, *http.Request) *APIError

func (fn deploymentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w,r); err != nil {
		http.Error(w, err.Message, err.Code)
	}
}
// ManageDeployment Creates, updates and deletes project deployments.
func deployment(w http.ResponseWriter, r *http.Request) *APIError{
	newDeployment := ProjectDeployment{}
	data, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(data, &newDeployment)

	switch method := r.Method; method {
	case "POST":
		return newDeployment.Insert(r.Context())
	case "DELETE":
		return newDeployment.Delete(r.Context())
	}
	return nil
}