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
	resp, err := k.Execute("get namespace")
	if err != nil {
		return nil, err
	}

	current, err := k.GetCurrentNamespace()
	if err != nil {
		return nil, err
	}

	outCh := resp.Readline()
	rawHeaders := <-outCh
	var namespaces []*Namespace
	for line := range outCh {
		ns := new(Namespace)
		if err := MakeResourceStruct(line, rawHeaders, ns); err != nil {
			return namespaces, err
		}

		if strings.EqualFold(ns.Name, current) {
			ns.Current = true
		}
		namespaces = append(namespaces, ns)
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
