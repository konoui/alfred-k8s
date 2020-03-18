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
	current, err := k.GetCurrentNamespace()
	if err != nil {
		return nil, err
	}

	// Note: NAME STATUS AGE
	arg := fmt.Sprintf("get namespace --no-headers")
	resp := k.Execute(arg)
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

	return namespaces, resp.err
}

// GetCurrentNamespace return current namespace name
func (k *Kubectl) GetCurrentNamespace() (string, error) {
	current, err := k.GetCurrentContext()
	if err != nil {
		return "", err
	}

	contexts, err := k.GetContexts()
	if err != nil {
		return "", err
	}

	for _, c := range contexts {
		if c.Name == current {
			return c.Namespace, nil
		}
	}

	return "", fmt.Errorf("found no namespace")
}

// SetNamespace configure namepsace
func (k *Kubectl) SetNamespace(context, ns string) error {
	arg := fmt.Sprintf("config set-context %s --namespace=%s", context, ns)
	resp := k.Execute(arg)
	return resp.err
}
