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
func (k *Kubectl) GetDeployments(all bool) ([]*Deployment, error) {
	if all {
		return k.getDeployments(allNamespaceFlag)
	}
	return k.getDeployments("")
}

func (k *Kubectl) getDeployments(ns string) ([]*Deployment, error) {
	// Note: NAME READY UP-TO-DATE AVAILABLE AGE
	// Note: NAMESPACE NAME READY UP-TO-DATE AVAILABLE AGE
	arg := fmt.Sprintf("get deployment %s --no-headers", ns)
	resp := k.Execute(arg)
	var deps []*Deployment
	for line := range resp.Readline() {
		rawData := strings.Fields(line)
		dep := generateDeployment(rawData)
		deps = append(deps, dep)
	}

	return deps, errors.Wrapf(resp.err, string(resp.stderr))
}

func generateDeployment(rawData []string) *Deployment {
	if len(rawData) == 5 {
		return &Deployment{
			Name:      rawData[0],
			Ready:     rawData[1],
			UpToDate:  rawData[2],
			Available: rawData[3],
			Age:       rawData[4],
		}
	}

	if len(rawData) == 6 {
		return &Deployment{
			Namespace: rawData[0],
			Name:      rawData[1],
			Ready:     rawData[2],
			UpToDate:  rawData[3],
			Available: rawData[4],
			Age:       rawData[5],
		}
	}

	msg := fmt.Sprintf("we assume that deployment information have 5 or 6 elements. but got %d elements. values: %v", len(rawData), rawData)
	panic(msg)
}
