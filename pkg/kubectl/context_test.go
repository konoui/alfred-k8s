package kubectl

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/konoui/alfred-k8s/pkg/executor"
)

var (
	testCurrentContext  = "test3-name"
	testContextsRawData = `*	test1-name  test1-cluster   test1-authinfo	<no value>
*	test2-name  test2-cluster   test2-authinfo	<no value>
*	test3-name  test3-cluster   test3-authinfo  test3-namespace`
)

var testContexts = []*Context{
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

var FakeContextFunc = func(args ...string) (*executor.Response, error) {
	if len(args) >= 2 {
		if args[1] == "current-context" {
			return &executor.Response{
				Stdout: []byte(testCurrentContext),
			}, nil
		}
		if args[1] == "view" {
			return &executor.Response{
				Stdout: []byte(testContextsRawData),
			}, nil
		}
	}
	return nil, fmt.Errorf("match no command args")
}

func TestGetCurrentContext(t *testing.T) {
	tests := []struct {
		name         string
		fakeExecutor executor.Executor
		want         string
	}{
		{
			name:         "get current context",
			fakeExecutor: NewFakeExecutor(FakeContextFunc),
			want:         testCurrentContext,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := SetupKubectl(t, tt.fakeExecutor)
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
		name         string
		fakeExecutor executor.Executor
		want         []*Context
	}{
		{
			name:         "view contexts",
			fakeExecutor: NewFakeExecutor(FakeContextFunc),
			want:         testContexts,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := SetupKubectl(t, tt.fakeExecutor)
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
		rawData        []string
		currentContext string
		want           *Context
	}{
		{
			name:           "generate current context",
			rawData:        []string{"CURRENT", "NAME", "CLUSTER", "AUTHINFO", "NAMESPACE"},
			currentContext: "NAME",
			want: &Context{
				Current:   true,
				Name:      "NAME",
				Namespace: "NAMESPACE",
			},
		},
		{
			name:           "generate not current context",
			rawData:        []string{"CURRENT", "NAME", "CLUSTER", "AUTHINFO", "NAMESPACE"},
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
			got := generateContext(tt.rawData, tt.currentContext)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}
