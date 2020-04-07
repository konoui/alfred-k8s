package kubectl

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/konoui/alfred-k8s/pkg/executor"
	"go.uber.org/goleak"
)

var testAllDeployments = []*Deployment{
	&Deployment{
		Namespace: "test1-namespace",
		Name:      "deployment-test1",
		Ready:     "1/1",
		UpToDate:  "1",
		Available: "1",
		Age:       "11d",
	},
	&Deployment{
		Namespace: "test2-namespace",
		Name:      "deployment-test2",
		Ready:     "2/2",
		UpToDate:  "2",
		Available: "1",
		Age:       "11d",
	},
}
var testDeployments = []*Deployment{
	&Deployment{
		Name:      "deployment-test1",
		Ready:     "1/1",
		UpToDate:  "1",
		Available: "1",
		Age:       "11d",
	},
	&Deployment{
		Name:      "deployment-test2",
		Ready:     "2/2",
		UpToDate:  "2",
		Available: "1",
		Age:       "11d",
	},
}

var FakeDeploymentFunc = func(t *testing.T, args ...string) (*executor.Response, error) {
	rawDataAllDeployments := GetByteFromTestFile(t, "testdata/raw-deployments-in-all-namespaces.txt")
	rawDataDeployments := GetByteFromTestFile(t, "testdata/raw-deployments.txt")
	if len(args) >= 4 {
		if args[1] == "deployment" && args[2] == allNamespaceFlag {
			return &executor.Response{
				Stdout: []byte(rawDataAllDeployments),
			}, nil
		}
	}
	if len(args) >= 2 {
		if args[1] == "deployment" {
			return &executor.Response{
				Stdout: []byte(rawDataDeployments),
			}, nil
		}
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

func TestGetDeployments(t *testing.T) {
	tests := []struct {
		name     string
		fakeFunc FakeFunc
		all      bool
		want     []*Deployment
	}{
		{
			name:     "list deployments",
			fakeFunc: FakeDeploymentFunc,
			want:     testDeployments,
		},
		{
			name:     "list deployments in all namespaces",
			fakeFunc: FakeDeploymentFunc,
			all:      true,
			want:     testAllDeployments,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, tt.fakeFunc)
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
