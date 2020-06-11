package kubectl

import (
	"fmt"
)

// Pod is kubectl get pod information
type Pod struct {
	Namespace string
	Name      string
	Ready     string
	Status    string
	Restarts  string
	Age       string
}

// GetPods return pods in current namespace
func (k *Kubectl) GetPods(all bool) ([]*Pod, error) {
	if all {
		return k.getPods(allNamespaceFlag)
	}
	return k.getPods("")
}

func (k *Kubectl) getPods(ns string) ([]*Pod, error) {
	// Note: NAME READY STATUS RESTARTS AGE
	// Note: NAMESPACE NAME READY STATUS RESTARTS AGE
	resp, err := k.Execute(fmt.Sprintf("get pod %s", ns))
	if err != nil {
		return nil, err
	}

	var pods []*Pod
	err = makeResourceStructSlice(resp, &pods)
	return pods, err
}
