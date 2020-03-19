package kubectl

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
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
func (k *Kubectl) GetDeployments() ([]*Deployment, error) {
	return k.getDeployments("")
}

// GetAllDeployments return all deployments in current context
func (k *Kubectl) GetAllDeployments() ([]*Deployment, error) {
	return k.getDeployments("--all-namespaces")
}

func (k *Kubectl) getDeployments(ns string) ([]*Deployment, error) {
	// Note: NAME READY UP-TO-DATE AVAILABLE AGE
	// Note: NAMESPACE NAME READY UP-TO-DATE AVAILABLE AGE
	arg := fmt.Sprintf("get deployment --no-headers %s", ns)
	resp := k.Execute(arg)
	var deps []*Deployment
	for line := range resp.Readline() {
		dInfo := strings.Fields(line)
		dep := generateDeployment(dInfo)
		deps = append(deps, dep)
	}

	return deps, errors.Wrapf(resp.err, string(resp.stderr))
}

func generateDeployment(dInfo []string) *Deployment {
	if len(dInfo) == 5 {
		return &Deployment{
			Name:      dInfo[0],
			Ready:     dInfo[1],
			UpToDate:  dInfo[2],
			Available: dInfo[3],
			Age:       dInfo[4],
		}
	}

	if len(dInfo) == 6 {
		return &Deployment{
			Namespace: dInfo[0],
			Name:      dInfo[1],
			Ready:     dInfo[2],
			UpToDate:  dInfo[3],
			Available: dInfo[4],
			Age:       dInfo[5],
		}
	}

	msg := fmt.Sprintf("we assume that deployment information have 5 or 6 elements. but got %d elements", len(dInfo))
	panic(msg)
}
