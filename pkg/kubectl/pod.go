package kubectl

import (
	"fmt"
	"strings"
)

// Pod presents `kubectl get pod -o wide`
type Pod struct {
	Name     string
	Ready    string
	Status   string
	Restarts string
	Age      string
	Node     string
}

// GetPods return pods in specific namespace
func GetPods(ns string) ([]*Pod, error) {
	return getPods(fmt.Sprintf("--namespace=%s", ns))
}

// GetAllPods return pods in all namespaces
func GetAllPods() ([]*Pod, error) {
	return getPods("--all-namespaces")
}

func getPods(ns string) ([]*Pod, error) {
	cmd := generateKubectlCmd(fmt.Sprintf("get pod %s", ns))
	resp := Execute(cmd)
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

	return pods, resp.err
}

func generatePod(podInfo, headers []string) *Pod {
	var pod Pod
	for i := range podInfo {
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
