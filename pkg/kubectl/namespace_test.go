package kubectl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/goleak"
)

var testNamespaces = []*Namespace{
	&Namespace{
		Name:   "test1-namepsace",
		Status: "Active",
		Age:    "11d",
	},
	&Namespace{
		Name:   "test2-namepsace",
		Status: "Active",
		Age:    "11d",
	},
	&Namespace{
		Current: true,
		Name:    "test3-namespace",
		Status:  "Active",
		Age:     "11d",
	},
	&Namespace{
		Name:   "default",
		Status: "Active",
		Age:    "11d",
	},
}

func TestGetCurrentNamespace(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "get current namespace",
			want: GetCurrentNamespaceFromtTestFile(t),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, nil)
			got, err := k.GetCurrentNamespace()
			if err != nil {
				t.Fatal(err)
			}

			if got != tt.want {
				t.Errorf("unexpected want: %v got: %v", tt.want, got)
			}
		})
	}
}

func TestGetNamespaces(t *testing.T) {
	tests := []struct {
		name string
		want []*Namespace
	}{
		{
			name: "list namespaces",
			want: testNamespaces,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, nil)
			got, err := k.GetNamespaces()
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}

func TestUseNamespace(t *testing.T) {
	tests := []struct {
		name string
		ns   string
	}{
		{
			name: "use namespace",
			ns:   "dummy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, nil)
			err := k.UseNamespace(tt.ns)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
