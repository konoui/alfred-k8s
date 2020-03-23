package cmd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{
			name: "override default value",
			want: &Config{
				Kubectl{
					Bin:         "/bin/ls",
					PluginPaths: []string{"/bin", "/usr/bin"},
				},
			},
		},
	}
	for _, tt := range tests {
		c, err := newConfig()
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(c, tt.want) {
			t.Errorf("want: \n%+v, got: \n%+v", tt.want, c)
		}
	}
}
