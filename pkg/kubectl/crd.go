package kubectl

import (
	"fmt"
	"strings"
)

const (
	crdTPL = `{{ range .items -}}
{{.metadata.name}}  {{.metadata.creationTimestamp}}
{{ end -}}`
	fieldName      = "NAME"
	fieldNamespace = "NAMESPACE"
	fieldAage      = "AGE"
)

// CRD Custom Resource Definition is kubectl get crd information
type CRD struct {
	Name      string
	CreatedAT string
}

// AnyResource is for other resources not supported
type AnyResource struct {
	Namespace string
	Name      string
	Age       string
}

// GetCRDs return crd in current namespace
func (k *Kubectl) GetCRDs() ([]*CRD, error) {
	arg := fmt.Sprintf("get crd -o go-template --template='%s'", crdTPL)
	resp := k.Execute(arg)
	var crds []*CRD
	for line := range resp.Readline() {
		rawData := strings.Fields(line)
		c := &CRD{
			Name:      rawData[0],
			CreatedAT: rawData[1],
		}
		crds = append(crds, c)
	}
	return crds, nil
}

// GetSpecificResources return specific resources in current namespace
func (k *Kubectl) GetSpecificResources(name string) ([]*AnyResource, error) {
	return k.getAnyResources(name, "")
}

// GetAllSpecificResources return specific resources in all namespaces
func (k *Kubectl) GetAllSpecificResources(name string) ([]*AnyResource, error) {
	return k.getAnyResources(name, allNamespaceFlag)
}

func (k *Kubectl) getAnyResources(name, ns string) ([]*AnyResource, error) {
	arg := fmt.Sprintf("get %s %s", name, ns)
	resp := k.Execute(arg)
	header := <-resp.Readline()
	headers := strings.Fields(header)

	var rs []*AnyResource
	for line := range resp.Readline() {
		rawData := strings.Fields(line)
		a := generateAnyResource(rawData, headers)
		// TODO
		if a.Name == fieldName {
			continue
		}
		rs = append(rs, a)
	}
	return rs, nil
}

func generateAnyResource(rawData, headers []string) *AnyResource {
	var c AnyResource
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
