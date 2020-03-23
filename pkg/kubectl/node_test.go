package kubectl

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/konoui/alfred-k8s/pkg/executor"
)

var (
	testNodesRawData = `node-1   Ready    <none>   11d   v1.15.10-eks-bac369
	node-2   Ready    <none>   11d   v1.15.10-eks-bac369`
)

var testNodes = []*Node{
	&Node{
		Name:    "node-1",
		Status:  "Ready",
		Roles:   "<none>",
		Age:     "11d",
		Version: "v1.15.10-eks-bac369",
	},
	&Node{
		Name:    "node-2",
		Status:  "Ready",
		Roles:   "<none>",
		Age:     "11d",
		Version: "v1.15.10-eks-bac369",
	},
}

var FakeNodeFunc = func(args ...string) (*executor.Response, error) {
	if len(args) >= 2 {
		if args[1] == "node" {
			return &executor.Response{
				Stdout: []byte(testNodesRawData),
			}, nil
		}
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

func TestGetNodes(t *testing.T) {
	tests := []struct {
		name         string
		fakeExecutor executor.Executor
		want         []*Node
	}{
		{
			name:         "list nodes",
			fakeExecutor: NewFakeExecutor(FakeNodeFunc),
			want:         testNodes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := SetupKubectl(t, tt.fakeExecutor)
			got, err := k.GetNodes()
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}
