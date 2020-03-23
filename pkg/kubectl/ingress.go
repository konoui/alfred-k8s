package kubectl

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Ingress is kubectl get ingress information
type Ingress struct {
	Namespace string
	Name      string
	Host      string
	Address   string
	Ports     string
	Age       string
}

// GetIngresses return ingresses in current namespace
func (k *Kubectl) GetIngresses() ([]*Ingress, error) {
	return k.getIngress("")
}

// GetAllIngresses return ingresses in all namespaces
func (k *Kubectl) GetAllIngresses() ([]*Ingress, error) {
	return k.getIngress(allNamespaceFlag)
}

func (k *Kubectl) getIngress(ns string) ([]*Ingress, error) {
	arg := fmt.Sprintf("get ingress %s --no-headers", ns)
	resp := k.Execute(arg)
	var ingresses []*Ingress
	for line := range resp.Readline() {
		rawData := strings.Fields(line)
		i := generateIngress(rawData)
		ingresses = append(ingresses, i)
	}
	return ingresses, errors.Wrapf(resp.err, string(resp.stderr))
}

func generateIngress(rawData []string) *Ingress {
	if len(rawData) == 5 {
		return &Ingress{
			Name:    rawData[0],
			Host:    rawData[1],
			Address: rawData[2],
			Ports:   rawData[3],
			Age:     rawData[4],
		}
	}
	if len(rawData) == 6 {
		return &Ingress{
			Namespace: rawData[0],
			Name:      rawData[1],
			Host:      rawData[2],
			Address:   rawData[3],
			Ports:     rawData[4],
			Age:       rawData[5],
		}
	}

	msg := fmt.Sprintf("we assume that ingress information have 5 or 6 elements. but got %d elements", len(rawData))
	panic(msg)
}
