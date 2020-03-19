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

// GetPods return pods in specific namespace
func (k *Kubectl) GetPods(ns string) ([]*Pod, error) {
	return k.getPods(fmt.Sprintf("--namespace=%s", ns))
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
		podInfo := strings.Fields(line)
		pod := generatePod(podInfo, words)
		// FIXME I don't understand why read repeated header from for-range.
		if pod.Name == "NAME" {
			continue
		}
		pods = append(pods, pod)
	}

	return pods, errors.Wrapf(resp.err, string(resp.stderr))
}

func generatePod(podInfo, headers []string) *Pod {
	var pod Pod
	for i := range podInfo {
		if strings.EqualFold(headers[i], "NAMESPACE") {
			pod.Namespace = podInfo[i]
		}
		if strings.EqualFold(headers[i], "NAME") {
			pod.Name = podInfo[i]
		}
		if strings.EqualFold(headers[i], "READY") {
			pod.Ready = podInfo[i]
		}
		if strings.EqualFold(headers[i], "STATUS") {
			pod.Status = podInfo[i]
		}
		if strings.EqualFold(headers[i], "RESTARTS") {
			pod.Restarts = podInfo[i]
		}
		if strings.EqualFold(headers[i], "AGE") {
			pod.Age = podInfo[i]
		}
	}

	return &pod
}
