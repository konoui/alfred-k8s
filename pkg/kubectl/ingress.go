package kubectl

import (
	"fmt"
)

// Ingress is kubectl get ingress information
type Ingress struct {
	Namespace string
	Name      string
	Hosts     string
	Address   string
	Ports     string
	Age       string
}

// GetIngresses return ingresses in current namespace
func (k *Kubectl) GetIngresses(all bool) ([]*Ingress, error) {
	if all {
		return k.getIngress(allNamespaceFlag)
	}
	return k.getIngress("")
}

func (k *Kubectl) getIngress(ns string) ([]*Ingress, error) {
	// Note: NAME	HOSTS	ADDRESS	PORTS	AGE
	// Note: NAMESPACE	NAME	HOSTS	ADDRESS	PORTS	AGE
	arg := fmt.Sprintf("get ingress %s", ns)
	resp, err := k.Execute(arg)
	if err != nil {
		return nil, err
	}

	var ingresses []*Ingress
	err = makeResourceStructSlice(resp, &ingresses)
	return ingresses, err
}
