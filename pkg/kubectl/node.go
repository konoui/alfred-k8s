package kubectl

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
	// Note: NAME STATUS ROLES AGE VERSION
	resp, err := k.Execute("get node")
	if err != nil {
		return nil, err
	}

	var nodes []*Node
	err = makeResourceStructSlice(resp, &nodes)
	return nodes, err
}
