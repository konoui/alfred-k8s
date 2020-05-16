package kubectl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/goleak"
)

var testAllDeployments = []*Deployment{
	{
		Namespace: "test1-namespace",
		Name:      "deployment-test1",
		Ready:     "1/1",
		UpToDate:  "1",
		Available: "1",
		Age:       "11d",
	},
	{
		Namespace: "test2-namespace",
		Name:      "deployment-test2",
		Ready:     "2/2",
		UpToDate:  "2",
		Available: "1",
		Age:       "11d",
	},
}
var testDeployments = []*Deployment{
	{
		Name:      "deployment-test1",
		Ready:     "1/1",
		UpToDate:  "1",
		Available: "1",
		Age:       "11d",
	},
	{
		Name:      "deployment-test2",
		Ready:     "2/2",
		UpToDate:  "2",
		Available: "1",
		Age:       "11d",
	},
}

func TestGetDeployments(t *testing.T) {
	tests := []struct {
		name string
		all  bool
		want []*Deployment
	}{
		{
			name: "list deployments",
			want: testDeployments,
		},
		{
			name: "list deployments in all namespaces",
			all:  true,
			want: testAllDeployments,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, nil)
			got, err := k.GetDeployments(tt.all)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}
