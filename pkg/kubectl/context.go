package kubectl

import (
	"context"
	"fmt"
)

// Context is kubectl get-contexts information
type Context struct {
	Current   bool
	Name      string
	Namespace string
}

// GetContexts return contexts
func (k *Kubectl) GetContexts() ([]*Context, error) {
	// Note: CURRENT NAME CLUSTER AUTHINFO NAMESPACE
	resp, err := k.Execute("config get-contexts")
	if err != nil {
		return nil, err
	}

	outCh := resp.Readline()
	rawHeaders := <-outCh
	var contexts []*Context
	for line := range outCh {
		l := new(struct{ Current string })
		c := new(Context)
		if err := MakeResourceStruct(line, rawHeaders, l); err != nil {
			return contexts, err
		}
		if err := MakeResourceStruct(line, rawHeaders, c); err != nil {
			return contexts, err
		}

		if l.Current == "*" {
			c.Current = true
		}
		contexts = append(contexts, c)
	}

	return contexts, nil
}

// GetCurrentContext return current configuration
func (k *Kubectl) GetCurrentContext() (string, error) {
	resp, err := k.Execute("config current-context")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return <-resp.ReadlineContext(ctx), err
}

// UseContext configure context
func (k *Kubectl) UseContext(c string) error {
	arg := fmt.Sprintf("config use-context %s", c)
	_, err := k.Execute(arg)
	return err
}
