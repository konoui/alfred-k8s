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

	outCh := resp.Readline()
	rawHeaders := <-outCh

	var pods []*Pod
	for line := range outCh {
		pod := new(Pod)
		if err := MakeResourceStruct(line, rawHeaders, pod); err != nil {
			return pods, err
		}
		pods = append(pods, pod)
	}

	return pods, nil
}
