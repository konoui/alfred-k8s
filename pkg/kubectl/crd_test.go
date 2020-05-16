package kubectl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/goleak"
)

var testCRDs = []*CRD{
	{
		Name:      "eniconfigs.crd.k8s.amazonaws.com",
		CreatedAT: "2020-03-11T12:23:16Z",
	},
	{
		Name:      "meshes.appmesh.k8s.aws",
		CreatedAT: "2020-03-11T12:32:15Z",
	},
}

func TestGetCRDs(t *testing.T) {
	tests := []struct {
		name string
		want []*CRD
	}{
		{
			name: "list CRDs",
			want: testCRDs,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, nil)
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
