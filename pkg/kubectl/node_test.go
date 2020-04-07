package kubectl

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/konoui/alfred-k8s/pkg/executor"
	"go.uber.org/goleak"
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

var FakeNodeFunc = func(t *testing.T, args ...string) (*executor.Response, error) {
	rawDataNodes := GetByteFromTestFile(t, "testdata/raw-nodes.txt")
	if len(args) >= 2 {
		if args[1] == "node" {
			return &executor.Response{
				Stdout: rawDataNodes,
			}, nil
		}
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

func TestGetNodes(t *testing.T) {
	tests := []struct {
		name     string
		fakeFunc FakeFunc
		want     []*Node
	}{
		{
			name:     "list nodes",
			fakeFunc: FakeNodeFunc,
			want:     testNodes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, tt.fakeFunc)
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
