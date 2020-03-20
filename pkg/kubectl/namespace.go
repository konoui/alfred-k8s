package kubectl

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
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
	resp := k.Execute(arg)

	current, err := k.GetCurrentNamespace()
	if err != nil {
		return nil, err
	}

	var namespaces []*Namespace
	for line := range resp.Readline() {
		nsInfo := strings.Fields(line)
		ns := Namespace{
			Name:   nsInfo[0],
			Status: nsInfo[1],
			Age:    nsInfo[2],
		}

		if nsInfo[0] == current {
			ns.Current = true
		}

		namespaces = append(namespaces, &ns)
	}

	return namespaces, errors.Wrapf(resp.err, string(resp.stderr))
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
	resp := k.Execute(arg)
	return errors.Wrapf(resp.err, string(resp.stderr))
}
