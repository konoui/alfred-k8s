package kubectl

import (
	"fmt"
	"strings"
)

// BaseResource is for other resources not supported
type BaseResource struct {
	Namespace string
	Name      string
	Age       string
}

// GetBaseResources return specific resources in current namespace
func (k *Kubectl) GetBaseResources(name string, all bool) ([]*BaseResource, error) {
	if all {
		return k.getBaseResources(name, allNamespaceFlag)
	}
	return k.getBaseResources(name, "")
}

func (k *Kubectl) getBaseResources(name, ns string) ([]*BaseResource, error) {
	arg := fmt.Sprintf("get %s %s", name, ns)
	resp, err := k.Execute(arg)
	stdout := resp.Readline()
	rawHeaders := <-stdout
	headers := strings.Fields(rawHeaders)

	var r []*BaseResource
	for line := range stdout {
		rawData := strings.Fields(line)
		a := generateBaseResource(rawData, headers)
		r = append(r, a)
	}
	return r, err
}

func generateBaseResource(rawData, headers []string) *BaseResource {
	var c BaseResource
	for i := range rawData {
		if headers[i] == knownNameField {
			c.Name = rawData[i]
		}
		if headers[i] == knownNamespaceField {
			c.Namespace = rawData[i]
		}
		if headers[i] == knownAageField {
			c.Age = rawData[i]
		}
	}
	return &c
}
