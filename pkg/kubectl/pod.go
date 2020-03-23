package kubectl

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
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
func (k *Kubectl) GetPods() ([]*Pod, error) {
	return k.getPods("")
}

// GetAllPods return pods in all namespaces
func (k *Kubectl) GetAllPods() ([]*Pod, error) {
	return k.getPods("--all-namespaces")
}

func (k *Kubectl) getPods(ns string) ([]*Pod, error) {
	// Note: NAME READY STATUS RESTARTS AGE
	// Note: NAMESPACE NAME READY STATUS RESTARTS AGE
	resp := k.Execute(fmt.Sprintf("get pod %s", ns))
	header := <-resp.Readline()
	words := strings.Fields(header)

	var pods []*Pod
	for line := range resp.Readline() {
		rawData := strings.Fields(line)
		pod := generatePod(rawData, words)
		// FIXME I don't understand why read repeated header from for-range.
		if pod.Name == "NAME" {
			continue
		}
		pods = append(pods, pod)
	}

	return pods, errors.Wrapf(resp.err, string(resp.stderr))
}

func generatePod(rawData, headers []string) *Pod {
	var pod Pod
	for i := range rawData {
		if strings.EqualFold(headers[i], "NAMESPACE") {
			pod.Namespace = rawData[i]
		}
		if strings.EqualFold(headers[i], "NAME") {
			pod.Name = rawData[i]
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
		if strings.EqualFold(headers[i], "AGE") {
			pod.Age = rawData[i]
		}
	}

	return &pod
}
