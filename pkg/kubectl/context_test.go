package kubectl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGenerateContext(t *testing.T) {
	tests := []struct {
		contextInfo    []string
		currentContext string
		want           *Context
	}{
		{
			contextInfo:    []string{"CURRENT", "NAME", "CLUSTER", "AUTHINFO", "NAMESPACE"},
			currentContext: "NAME",
			want: &Context{
				Current:   true,
				Name:      "NAME",
				Namespace: "NAMESPACE",
			},
		},
		{
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
		got := generateContext(tt.contextInfo, tt.currentContext)
		if diff := cmp.Diff(got, tt.want); diff != "" {
			t.Errorf("+want -got\n%+v", diff)
		}
	}
}
