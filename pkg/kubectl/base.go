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

const (
	fieldName      = "NAME"
	fieldNamespace = "NAMESPACE"
	fieldAage      = "AGE"
)

// GetBaseResources return specific resources in current namespace
func (k *Kubectl) GetBaseResources(name string, all bool) ([]*BaseResource, error) {
	if all {
		return k.getBaseResources(name, allNamespaceFlag)
	}
	return k.getBaseResources(name, "")
}

func (k *Kubectl) getBaseResources(name, ns string) ([]*BaseResource, error) {
	arg := fmt.Sprintf("get %s %s", name, ns)
	resp := k.Execute(arg)
	header := <-resp.Readline()
	headers := strings.Fields(header)

	var rs []*BaseResource
	for line := range resp.Readline() {
		rawData := strings.Fields(line)
		a := generateBaseResource(rawData, headers)
		// TODO
		if a.Name == fieldName {
			continue
		}
		rs = append(rs, a)
	}
	return rs, nil
}

func generateBaseResource(rawData, headers []string) *BaseResource {
	var c BaseResource
	for i := range rawData {
		if headers[i] == fieldName {
			c.Name = rawData[i]
		}
		if headers[i] == fieldNamespace {
			c.Namespace = rawData[i]
		}
		if headers[i] == fieldAage {
			c.Age = rawData[i]
		}
	}
	return &c
}
