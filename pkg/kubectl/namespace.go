package kubectl

import (
	"fmt"
	"strings"
)

// Namespace is kubectl get ns information
type Namespace struct {
	Name    string
	Current bool
	Status  string
	Age     string
}

// GetNamespaces return namespaces in current context
func (k *Kubectl) GetNamespaces() ([]*Namespace, error) {
	// Note: NAME STATUS AGE
	arg := fmt.Sprintf("get namespace --no-headers")
	resp, err := k.Execute(arg)
	if err != nil {
		return nil, err
	}

	current, err := k.GetCurrentNamespace()
	if err != nil {
		return nil, err
	}

	var namespaces []*Namespace
	for line := range resp.Readline() {
		rawData := strings.Fields(line)
		ns := Namespace{
			Name:   rawData[0],
			Status: rawData[1],
			Age:    rawData[2],
		}

		if rawData[0] == current {
			ns.Current = true
		}

		namespaces = append(namespaces, &ns)
	}

	return namespaces, nil
}

// GetCurrentNamespace return current namespace name
func (k *Kubectl) GetCurrentNamespace() (string, error) {
	contexts, err := k.GetContexts()
	if err != nil {
		return "", err
	}

	for _, c := range contexts {
		if c.Current {
			return c.Namespace, nil
		}
	}

	// namespace will be empty if namespace does not set in kubeconfig
	return "", nil
}

// UseNamespace configure namepsace in current context
func (k *Kubectl) UseNamespace(ns string) error {
	context, err := k.GetCurrentContext()
	if err != nil {
		return err
	}

	arg := fmt.Sprintf("config set-context %s --namespace=%s", context, ns)
	_, err = k.Execute(arg)
	return err
}
