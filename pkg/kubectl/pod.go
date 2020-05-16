package kubectl

import (
	"fmt"
	"strings"
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
	stdout := resp.Readline()
	rawHeaders := <-stdout
	headers := strings.Fields(rawHeaders)

	var pods []*Pod
	for line := range stdout {
		rawData := strings.Fields(line)
		pod := generatePod(rawData, headers)
		pods = append(pods, pod)
	}

	return pods, err
}

func generatePod(rawData, headers []string) *Pod {
	var pod Pod
	for i := range rawData {
		if strings.EqualFold(headers[i], knownNamespaceField) {
			pod.Namespace = rawData[i]
		}
		if strings.EqualFold(headers[i], knownNameField) {
			pod.Name = rawData[i]
		}
		if strings.EqualFold(headers[i], knownAageField) {
			pod.Age = rawData[i]
		}
		if strings.EqualFold(headers[i], "READY") {
			pod.Ready = rawData[i]
		}
		if strings.EqualFold(headers[i], "STATUS") {
			pod.Status = rawData[i]
		}
		if strings.EqualFold(headers[i], "RESTARTS") {
			pod.Restarts = rawData[i]
		}
	}

	return &pod
}
