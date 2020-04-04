package kubectl

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/konoui/alfred-k8s/pkg/executor"
	"go.uber.org/goleak"
)

var testCRDsRawData = `eniconfigs.crd.k8s.amazonaws.com   2020-03-11T12:23:16Z
meshes.appmesh.k8s.aws             2020-03-11T12:32:15Z`

var testCRDs = []*CRD{
	&CRD{
		Name:      "eniconfigs.crd.k8s.amazonaws.com",
		CreatedAT: "2020-03-11T12:23:16Z",
	},
	&CRD{
		Name:      "meshes.appmesh.k8s.aws",
		CreatedAT: "2020-03-11T12:32:15Z",
	},
}

var FakeCRDFunc = func(args ...string) (*executor.Response, error) {
	if len(args) >= 1 {
		if args[1] == "crd" {
			return &executor.Response{
				Stdout: []byte(testCRDsRawData),
			}, nil
		}
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

func TestGetCRDs(t *testing.T) {
	tests := []struct {
		name         string
		fakeExecutor executor.Executor
		want         []*CRD
	}{
		{
			name:         "list CRDs",
			fakeExecutor: NewFakeExecutor(FakeCRDFunc),
			want:         testCRDs,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, tt.fakeExecutor)
			got, err := k.GetCRDs()
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}
