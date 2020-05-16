package kubectl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/goleak"
)

var testNodes = []*Node{
	{
		Name:    "node-1",
		Status:  "Ready",
		Roles:   "<none>",
		Age:     "11d",
		Version: "v1.15.10-eks-bac369",
	},
	{
		Name:    "node-2",
		Status:  "Ready",
		Roles:   "<none>",
		Age:     "11d",
		Version: "v1.15.10-eks-bac369",
	},
}

func TestGetNodes(t *testing.T) {
	tests := []struct {
		name string
		want []*Node
	}{
		{
			name: "list nodes",
			want: testNodes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, nil)
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
