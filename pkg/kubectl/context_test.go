package kubectl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/goleak"
)

var testContexts = []*Context{
	{
		Name:      "test1-name",
		Namespace: "",
	},
	{
		Name:      "test2-name",
		Namespace: "",
	},
	{
		Current:   true,
		Name:      "test3-name",
		Namespace: "test3-namespace",
	},
}

func TestGetCurrentContext(t *testing.T) {
	tests := []struct {
		name     string
		fakeFunc FakeFunc
		want     string
	}{
		{
			name: "get current context",
			want: GetCurrentContextFromtTestFile(t),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, nil)
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
		name string
		want []*Context
	}{
		{
			name: "view contexts",
			want: testContexts,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, nil)
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

func TestUseContext(t *testing.T) {
	tests := []struct {
		name    string
		context string
	}{
		{
			name:    "use context",
			context: "dummy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, nil)
			err := k.UseContext(tt.context)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
