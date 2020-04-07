package kubectl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/goleak"
)

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

func TestGetCRDs(t *testing.T) {
	tests := []struct {
		name     string
		fakeFunc FakeFunc
		want     []*CRD
	}{
		{
			name:     "list CRDs",
			fakeFunc: FakeCRDFunc,
			want:     testCRDs,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, tt.fakeFunc)
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
