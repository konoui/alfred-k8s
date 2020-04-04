package kubectl

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
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
	stdout := resp.Readline()
	rawHeaders := <-stdout
	headers := strings.Fields(rawHeaders)

	var rs []*BaseResource
	for line := range stdout {
		rawData := strings.Fields(line)
		a := generateBaseResource(rawData, headers)
		rs = append(rs, a)
	}
	return rs, errors.Wrapf(resp.err, string(resp.stderr))
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
