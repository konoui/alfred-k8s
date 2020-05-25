package kubectl

import (
	"fmt"
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
	if err != nil {
		return nil, err
	}

	outCh := resp.Readline()
	rawHeaders := <-outCh

	var r []*BaseResource
	for line := range outCh {
		a := new(BaseResource)
		if err := MakeResourceStruct(line, rawHeaders, a); err != nil {
			return r, err
		}
		r = append(r, a)
	}

	return r, nil
}
