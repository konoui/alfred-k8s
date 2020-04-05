package kubectl

import (
	"fmt"
	"strings"
)

const crdTPL = `{{ range .items -}}
{{.metadata.name}}  {{.metadata.creationTimestamp}}
{{ end -}}`

// CRD Custom Resource Definition is kubectl get crd information
type CRD struct {
	Name      string
	CreatedAT string
}

// GetCRDs return crd in current namespace
func (k *Kubectl) GetCRDs() ([]*CRD, error) {
	arg := fmt.Sprintf("get crd -o go-template --template='%s'", crdTPL)
	resp, err := k.Execute(arg)
	var crds []*CRD
	for line := range resp.Readline() {
		rawData := strings.Fields(line)
		c := &CRD{
			Name:      rawData[0],
			CreatedAT: rawData[1],
		}
		crds = append(crds, c)
	}
	return crds, err
}
