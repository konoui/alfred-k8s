package kubectl

import (
	"strings"
)

// Node is kubectl get node information
type Node struct {
	Name    string
	Status  string
	Roles   string
	Age     string
	Version string
}

// GetNodes return nodes
func (k *Kubectl) GetNodes() ([]*Node, error) {
	resp, err := k.Execute("get node --no-headers")

	var nodes []*Node
	for line := range resp.Readline() {
		rawData := strings.Fields(line)
		n := &Node{
			Name:    rawData[0],
			Status:  rawData[1],
			Roles:   rawData[2],
			Age:     rawData[3],
			Version: rawData[4],
		}
		nodes = append(nodes, n)
	}

	return nodes, err
}
