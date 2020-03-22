package kubectl

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/konoui/alfred-k8s/pkg/executor"
)

type fakeNamespaceExecutor struct {
	namespaces      *executor.Response
	namespacesErr   error
	contextExecutor executor.Executor
}

func (e *fakeNamespaceExecutor) Exec(args ...string) (*executor.Response, error) {
	if len(args) >= 2 {
		if args[1] == "namespace" {
			return e.namespaces, e.namespacesErr
		}
		// Note: get current namespace  and namespaces call context function
		return e.contextExecutor.Exec(args...)

	}
	return nil, fmt.Errorf("match no command args")
}

const (
	testCurrentNamespace = "test3-namespace"
	testNamespaceInputs  = `test1-namepsace	Active	11d
	test2-namepsace	Active	11d
	test3-namespace	Active	11d
	default	Active	11d`
)

var testNamespaceOutputs = []*Namespace{
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

var testNamespaceExecutor = &fakeNamespaceExecutor{
	namespaces: &executor.Response{
		Stdout: []byte(testNamespaceInputs),
	},
	contextExecutor: testContextExecutor,
}

func TestGetCurrentNamespace(t *testing.T) {
	tests := []struct {
		name         string
		fakeExecutor executor.Executor
		want         string
	}{
		{
			name:         "get current namespace",
			fakeExecutor: testContextExecutor,
			want:         testCurrentNamespace,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := testSetupKubectl(t, tt.fakeExecutor)
			got, err := k.GetCurrentNamespace()
			if err != nil {
				t.Fatal(err)
			}

			if got != tt.want {
				t.Errorf("unexpected want: %v\ngot: %v", tt.want, got)
			}
		})
	}
}

func TestGetNamespaces(t *testing.T) {
	tests := []struct {
		name         string
		fakeExecutor executor.Executor
		want         []*Namespace
	}{
		{
			name:         "list namespaces",
			fakeExecutor: testNamespaceExecutor,
			want:         testNamespaceOutputs,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := testSetupKubectl(t, tt.fakeExecutor)
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
