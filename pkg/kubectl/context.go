package kubectl

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const (
	// see https://golang.org/pkg/text/template/#Template.Option
	noValue    = "<no value>"
	dummyValue = "DUMMY"
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
	// Note: CURRENT is dummy VALUE
	resp := k.Execute("config view -o go-template-file=context.tpl")
	current, err := k.GetCurrentContext()
	if err != nil {
		return nil, err
	}

	var contexts []*Context
	for line := range resp.Readline() {
		rawData := strings.Fields(strings.Replace(line, noValue, dummyValue, -1))
		c := generateContext(rawData, current)
		contexts = append(contexts, c)
	}

	return contexts, errors.Wrapf(resp.err, string(resp.stderr))
}

func generateContext(rawData []string, current string) *Context {
	if len(rawData) != 5 {
		msg := fmt.Sprintf("we assume that context information have 5 elements. but got %d. values: %v", len(rawData), rawData)
		panic(msg)
	}

	for i := range rawData {
		if rawData[i] == dummyValue {
			rawData[i] = ""
		}
	}

	c := Context{
		Current:   false,
		Name:      rawData[1],
		Namespace: rawData[4],
	}
	if c.Name == current {
		c.Current = true
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

// UseContext configure context
func (k *Kubectl) UseContext(c string) error {
	arg := fmt.Sprintf("config use-context %s", c)
	resp := k.Execute(arg)
	return errors.Wrapf(resp.err, string(resp.stderr))
}
