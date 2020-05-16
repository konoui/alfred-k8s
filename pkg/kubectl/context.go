package kubectl

import (
	"context"
	"fmt"
	"strings"
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
	indexMap := makeIndexMap(rawHeaders)

	var contexts []*Context
	for line := range outCh {
		c := makeContext(line, indexMap)
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

func makeIndexMap(rawHeaders string) (indexMap map[string]int) {
	indexMap = make(map[string]int)
	headers := strings.Fields(rawHeaders)
	for _, h := range headers {
		indexMap[h] = strings.Index(rawHeaders, h)
	}
	return
}

func makeContext(line string, indexMap map[string]int) *Context {
	var c Context
	for key, start := range indexMap {
		value := ""
		if len(line) > start {
			value = strings.Fields(line[start:])[0]
		}

		if strings.EqualFold(key, knownNamespaceField) {
			c.Namespace = value
		}
		if strings.EqualFold(key, knownNameField) {
			c.Name = value
		}
		if strings.EqualFold(key, "CURRENT") {
			if value == "*" {
				c.Current = true
			}
		}
	}
	return &c
}
