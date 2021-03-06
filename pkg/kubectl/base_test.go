package kubectl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/goleak"
)

var testBaseAllPods = []*BaseResource{
	{
		Namespace: "test1-namespace",
		Name:      "test1-pod",
		Age:       "11d",
	},
	{
		Namespace: "test2-namespace",
		Name:      "test2-pod",
		Age:       "11d",
	},
}

var testBasePods = []*BaseResource{
	{
		Name: "test1-pod",
		Age:  "11d",
	},
	{
		Name: "test2-pod",
		Age:  "11d",
	},
}

func TestPodBaseResource(t *testing.T) {
	tests := []struct {
		name string
		all  bool
		want []*BaseResource
	}{
		{
			name: "list pods for base resource",
			want: testBasePods,
		},
		{
			name: "list pods in all namespaces for base resource",
			all:  true,
			want: testBaseAllPods,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, nil)
			got, err := k.GetBaseResources("pod", tt.all)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}
