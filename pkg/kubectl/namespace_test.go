package kubectl

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/konoui/alfred-k8s/pkg/executor"
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

func getCurrentNamespace(t *testing.T) string {
	c := GetStringFromTestFile(t, "testdata/raw-current-namespace.txt")
	return strings.Replace(c, "\n", "", -1)
}

var FakeNamespaceFunc = func(t *testing.T, args ...string) (*executor.Response, error) {
	rawDataNamespaces := GetByteFromTestFile(t, "testdata/raw-namespaces.txt")
	if len(args) >= 2 {
		if args[1] == "namespace" {
			return &executor.Response{
				Stdout: rawDataNamespaces,
			}, nil
		}
		// Note: get current namespace and namespaces call context function
		return FakeContextFunc(t, args...)
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

func TestGetCurrentNamespace(t *testing.T) {
	tests := []struct {
		name     string
		fakeFunc FakeFunc
		want     string
	}{
		{
			name:     "get current namespace",
			fakeFunc: FakeNamespaceFunc,
			want:     getCurrentNamespace(t),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, tt.fakeFunc)
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
		name     string
		fakeFunc FakeFunc
		want     []*Namespace
	}{
		{
			name:     "list namespaces",
			fakeFunc: FakeNamespaceFunc,
			want:     testNamespaces,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, tt.fakeFunc)
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
