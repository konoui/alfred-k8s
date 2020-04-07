package kubectl

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/konoui/alfred-k8s/pkg/executor"
	"go.uber.org/goleak"
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

func getCurrentContext(t *testing.T) string {
	c := GetStringFromTestFile(t, "testdata/raw-current-context.txt")
	return strings.Replace(c, "\n", "", -1)
}

var FakeContextFunc = func(t *testing.T, args ...string) (*executor.Response, error) {
	rawDataCurrentContext := GetByteFromTestFile(t, "testdata/raw-current-context.txt")
	rawDataContexts := GetByteFromTestFile(t, "testdata/raw-contexts.txt")
	if len(args) >= 2 {
		if args[1] == "current-context" {
			return &executor.Response{
				Stdout: rawDataCurrentContext,
			}, nil
		}
		if args[1] == "view" {
			return &executor.Response{
				Stdout: rawDataContexts,
			}, nil
		}
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

func TestGetCurrentContext(t *testing.T) {
	tests := []struct {
		name     string
		fakeFunc FakeFunc
		want     string
	}{
		{
			name:     "get current context",
			fakeFunc: FakeContextFunc,
			want:     getCurrentContext(t),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, tt.fakeFunc)
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
		name     string
		fakeFunc FakeFunc
		want     []*Context
	}{
		{
			name:     "view contexts",
			fakeFunc: FakeContextFunc,
			want:     testContexts,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, tt.fakeFunc)
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
