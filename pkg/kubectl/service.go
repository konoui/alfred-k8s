package kubectl

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Service is kubectl get service information
type Service struct {
	Namespace  string
	Name       string
	Type       string
	ClusterIP  string
	ExternalIP string
	Ports      string
	Age        string
}

// GetServices return services in current namespace
func (k *Kubectl) GetServices() ([]*Service, error) {
	return k.getServices("")
}

// GetAllServices return services in current context
func (k *Kubectl) GetAllServices() ([]*Service, error) {
	return k.getServices("--all-namespaces")
}

func (k *Kubectl) getServices(ns string) ([]*Service, error) {
	// Note: NAME TYPE CLUSTER-IP EXTERNAL-IP PORT(S) AGE
	// Note: NAMESPACE NAME TYPE CLUSTER-IP EXTERNAL-IP PORT(S) AGE
	arg := fmt.Sprintf("get service --no-headers %s", ns)
	resp := k.Execute(arg)
	var svcs []*Service
	for line := range resp.Readline() {
		rawData := strings.Fields(line)
		svc := generateService(rawData)
		svcs = append(svcs, svc)
	}

	return svcs, errors.Wrapf(resp.err, string(resp.stderr))
}

func generateService(rawData []string) *Service {
	if len(rawData) == 6 {
		return &Service{
			Name:       rawData[0],
			Type:       rawData[1],
			ClusterIP:  rawData[2],
			ExternalIP: rawData[3],
			Ports:      rawData[4],
			Age:        rawData[5],
		}
	}

	if len(rawData) == 7 {
		return &Service{
			Namespace:  rawData[0],
			Name:       rawData[1],
			Type:       rawData[2],
			ClusterIP:  rawData[3],
			ExternalIP: rawData[4],
			Ports:      rawData[5],
			Age:        rawData[6],
		}
	}

	msg := fmt.Sprintf("we assume that service information have 6 or 7 elements. but got %d elements", len(rawData))
	panic(msg)
}
