package kubectl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/goleak"
)

var testAllIngresses = []*Ingress{
	&Ingress{
		Namespace: "test1-namespace",
		Name:      "test1-ingress",
		Hosts:     "*",
		Address:   "ingress1.hoge.hoge",
		Ports:     "80",
		Age:       "24h",
	},
	&Ingress{
		Namespace: "test2-namespace",
		Name:      "test2-ingress",
		Hosts:     "*",
		Address:   "ingress2.hoge.hoge",
		Ports:     "80",
		Age:       "24h",
	},
}
var testIngresses = []*Ingress{
	&Ingress{
		Name:    "test1-ingress",
		Hosts:   "*",
		Address: "ingress1.hoge.hoge",
		Ports:   "80",
		Age:     "24h",
	},
	&Ingress{
		Name:    "test2-ingress",
		Hosts:   "*",
		Address: "ingress2.hoge.hoge",
		Ports:   "80",
		Age:     "24h",
	},
}

func TestGetIngresses(t *testing.T) {
	tests := []struct {
		name     string
		fakeFunc FakeFunc
		all      bool
		want     []*Ingress
	}{
		{
			name:     "list ingresses",
			fakeFunc: FakeIngressFunc,
			want:     testIngresses,
		},
		{
			name:     "list deployments in all namespaces",
			fakeFunc: FakeIngressFunc,
			all:      true,
			want:     testAllIngresses,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, tt.fakeFunc)
			got, err := k.GetIngresses(tt.all)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}
