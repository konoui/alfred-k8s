package kubectl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/goleak"
)

var testAllIngresses = []*Ingress{
	{
		Namespace: "test1-namespace",
		Name:      "test1-ingress",
		Hosts:     "*",
		Address:   "ingress1.hoge.hoge",
		Ports:     "80",
		Age:       "24h",
	},
	{
		Namespace: "test2-namespace",
		Name:      "test2-ingress",
		Hosts:     "*",
		Address:   "ingress2.hoge.hoge",
		Ports:     "80",
		Age:       "24h",
	},
}
var testIngresses = []*Ingress{
	{
		Name:    "test1-ingress",
		Hosts:   "*",
		Address: "ingress1.hoge.hoge",
		Ports:   "80",
		Age:     "24h",
	},
	{
		Name:    "test2-ingress",
		Hosts:   "*",
		Address: "ingress2.hoge.hoge",
		Ports:   "80",
		Age:     "24h",
	},
}

func TestGetIngresses(t *testing.T) {
	tests := []struct {
		name string
		all  bool
		want []*Ingress
	}{
		{
			name: "list ingresses",
			want: testIngresses,
		},
		{
			name: "list deployments in all namespaces",
			all:  true,
			want: testAllIngresses,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, nil)
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
