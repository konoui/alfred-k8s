package kubectl

import (
	"fmt"
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
func (k *Kubectl) GetServices(all bool) ([]*Service, error) {
	if all {
		return k.getServices(allNamespaceFlag)
	}
	return k.getServices("")
}

func (k *Kubectl) getServices(ns string) ([]*Service, error) {
	// Note: NAME TYPE CLUSTER-IP EXTERNAL-IP PORT(S) AGE
	// Note: NAMESPACE NAME TYPE CLUSTER-IP EXTERNAL-IP PORT(S) AGE
	arg := fmt.Sprintf("get service %s", ns)
	resp, err := k.Execute(arg)
	if err != nil {
		return nil, err
	}

	outCh := resp.Readline()
	rawHeaders := <-outCh

	var svcs []*Service
	for line := range outCh {
		svc := new(Service)
		if err := MakeResourceStruct(line, rawHeaders, svc); err != nil {
			return svcs, err
		}
		svcs = append(svcs, svc)
	}

	return svcs, nil
}
