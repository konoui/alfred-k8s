package kubectl

import (
	"fmt"
	"strings"
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
	arg := fmt.Sprintf("get ingress %s --no-headers", ns)
	resp, err := k.Execute(arg)
	var ingresses []*Ingress
	for line := range resp.Readline() {
		rawData := strings.Fields(line)
		i := generateIngress(rawData)
		ingresses = append(ingresses, i)
	}
	return ingresses, err
}

func generateIngress(rawData []string) *Ingress {
	if len(rawData) == 5 {
		return &Ingress{
			Name:    rawData[0],
			Hosts:   rawData[1],
			Address: rawData[2],
			Ports:   rawData[3],
			Age:     rawData[4],
		}
	}
	if len(rawData) == 6 {
		return &Ingress{
			Namespace: rawData[0],
			Name:      rawData[1],
			Hosts:     rawData[2],
			Address:   rawData[3],
			Ports:     rawData[4],
			Age:       rawData[5],
		}
	}

	msg := fmt.Sprintf("we assume that ingress information have 5 or 6 elements. but got %d elements. values: %v", len(rawData), rawData)
	panic(msg)
}
