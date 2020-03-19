package kubectl

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
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
	resp := k.Execute("config get-contexts --no-headers")
	var contexts []*Context
	for line := range resp.Readline() {
		cInfo := strings.Fields(line)
		c := generateContext(cInfo)
		contexts = append(contexts, c)
	}

	return contexts, errors.Wrapf(resp.err, string(resp.stderr))
}

func generateContext(cInfo []string) *Context {
	var c Context
	if len(cInfo) < 3 {
		panic("we assume that context information have NAME, CLUSTER and AUTHINFO elements at least.")
	}
	currentMarker := cInfo[0]
	if currentMarker == "*" {
		// current context case
		c.Current = true
		c.Name = cInfo[1]
		if len(cInfo) < 5 {
			c.Namespace = "default"
		} else {
			c.Namespace = cInfo[4]
		}
		return &c
	}

	// if not current context case, 0 element will be context name
	c.Name = cInfo[0]
	if len(cInfo) < 4 {
		c.Namespace = "default"
	} else {
		c.Namespace = cInfo[3]
	}
	return &c
}

// GetCurrentContext return current configuration
func (k *Kubectl) GetCurrentContext() (string, error) {
	resp := k.Execute("config current-context")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return <-resp.ReadlineContext(ctx), errors.Wrapf(resp.err, string(resp.stderr))
}

// SetContext configure context
func (k *Kubectl) SetContext(c string) error {
	arg := fmt.Sprintf("config use-context %s", c)
	resp := k.Execute(arg)
	return errors.Wrapf(resp.err, string(resp.stderr))
}
