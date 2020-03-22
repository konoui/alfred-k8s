package kubectl

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/konoui/alfred-k8s/pkg/executor"
)

// type fakeExecutor struct {
// 	response *executor.Response
// 	err      error
// }

// func (e *fakeExecutor) Exec(args ...string) (*executor.Response, error) {
// 	return e.response, e.err
// }

type fakeContextExecutor struct {
	current     *executor.Response
	currentErr  error
	contexts    *executor.Response
	contextsErr error
}

func (e *fakeContextExecutor) Exec(args ...string) (*executor.Response, error) {
	if len(args) >= 2 {
		if args[1] == "current-context" {
			return e.current, e.currentErr
		}
		if args[1] == "view" {
			return e.contexts, e.contextsErr
		}
	}
	return nil, fmt.Errorf("match no command args")
}

const (
	testCurrentContext = "test3-name"
	testContextInputs  = `*	test1-name  test1-cluster   test1-authinfo	<no value>
*	test2-name  test2-cluster   test2-authinfo	<no value>
*	test3-name  test3-cluster   test3-authinfo  test3-namespace`
)

var testContextOutputs = []*Context{
	&Context{
		Name:      "test1-name",
		Namespace: "",
	},
	&Context{
		Name:      "test2-name",
		Namespace: "",
	},
	&Context{
		Current:   true,
		Name:      "test3-name",
		Namespace: "test3-namespace",
	},
}

var testContextExecutor = &fakeContextExecutor{
	current: &executor.Response{
		Stdout: []byte(testCurrentContext),
	},
	contexts: &executor.Response{
		Stdout: []byte(testContextInputs),
	},
}

func testSetupKubectl(t *testing.T, fake executor.Executor) *Kubectl {
	t.Helper()
	k, err := New(OptionExecutor(fake))
	if err != nil {
		t.Fatal(err)
	}
	return k
}

func TestGetCurrentContext(t *testing.T) {
	tests := []struct {
		name         string
		fakeExecutor executor.Executor
		want         string
	}{
		{
			name:         "get current context",
			fakeExecutor: testContextExecutor,
			want:         testCurrentContext,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := testSetupKubectl(t, tt.fakeExecutor)

			got, err := k.GetCurrentContext()
			if err != nil {
				t.Fatal(err)
			}

			if got != tt.want {
				t.Errorf("unexpected want: %v\ngot: %v", tt.want, got)
			}
		})
	}
}

func TestGetContexts(t *testing.T) {
	tests := []struct {
		name                string
		fakeContextExecutor executor.Executor
		want                []*Context
	}{
		{
			name:                "view contexts",
			fakeContextExecutor: testContextExecutor,
			want:                testContextOutputs,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := testSetupKubectl(t, tt.fakeContextExecutor)

			got, err := k.GetContexts()
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}

func TestGenerateContext(t *testing.T) {
	tests := []struct {
		name           string
		contextInfo    []string
		currentContext string
		want           *Context
	}{
		{
			name:           "generate current context",
			contextInfo:    []string{"CURRENT", "NAME", "CLUSTER", "AUTHINFO", "NAMESPACE"},
			currentContext: "NAME",
			want: &Context{
				Current:   true,
				Name:      "NAME",
				Namespace: "NAMESPACE",
			},
		},
		{
			name:           "generate not current context",
			contextInfo:    []string{"CURRENT", "NAME", "CLUSTER", "AUTHINFO", "NAMESPACE"},
			currentContext: "NOT CURRENT CONTEXT",
			want: &Context{
				Current:   false,
				Name:      "NAME",
				Namespace: "NAMESPACE",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateContext(tt.contextInfo, tt.currentContext)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}
