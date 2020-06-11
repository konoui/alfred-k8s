package kubectl

import (
	"fmt"
)

// Deployment is kubectl get deployment information
type Deployment struct {
	Namespace string
	Name      string
	Ready     string
	UpToDate  string
	Available string
	Age       string
}

// GetDeployments return deployments in current namespace
func (k *Kubectl) GetDeployments(all bool) ([]*Deployment, error) {
	if all {
		return k.getDeployments(allNamespaceFlag)
	}
	return k.getDeployments("")
}

func (k *Kubectl) getDeployments(ns string) ([]*Deployment, error) {
	// Note: NAME READY UP-TO-DATE AVAILABLE AGE
	// Note: NAMESPACE NAME READY UP-TO-DATE AVAILABLE AGE
	arg := fmt.Sprintf("get deployment %s", ns)
	resp, err := k.Execute(arg)
	if err != nil {
		return nil, err
	}

	var deps []*Deployment
	err = makeResourceStructSlice(resp, &deps)
	return deps, err
}
