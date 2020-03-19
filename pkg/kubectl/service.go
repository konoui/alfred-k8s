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
		sInfo := strings.Fields(line)
		svc := generateService(sInfo)
		svcs = append(svcs, svc)
	}

	return svcs, errors.Wrapf(resp.err, string(resp.stderr))
}

func generateService(sInfo []string) *Service {
	if len(sInfo) == 6 {
		return &Service{
			Name:       sInfo[0],
			Type:       sInfo[1],
			ClusterIP:  sInfo[2],
			ExternalIP: sInfo[3],
			Ports:      sInfo[4],
			Age:        sInfo[5],
		}
	}

	if len(sInfo) == 7 {
		return &Service{
			Namespace:  sInfo[0],
			Name:       sInfo[1],
			Type:       sInfo[2],
			ClusterIP:  sInfo[3],
			ExternalIP: sInfo[4],
			Ports:      sInfo[5],
			Age:        sInfo[6],
		}
	}

	msg := fmt.Sprintf("we assume that service information have 6 or 7 elements. but got %d elements", len(sInfo))
	panic(msg)
}
