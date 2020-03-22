package kubectl

import (
	"reflect"
	"testing"
)

func TestNewKubectl(t *testing.T) {
	tests := []struct {
		name      string
		options   []Option
		want      *Kubectl
		expectErr bool
	}{
		{
			name: "default value",
			want: &Kubectl{
				cmd:        newCommand("/usr/local/bin/kubectl"),
				pluginPath: "/usr/local/bin/",
			},
		},
		{
			name: "options value",
			options: []Option{
				OptionBinary(knownBinary),
				OptionPluginPath(knownBinPath),
			},
			want: &Kubectl{
				cmd:        newCommand(knownBinary),
				pluginPath: knownBinPath,
			},
			// TODO unexptected bin path case.
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.options...)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("want: %+v\ngot: %+v", tt.want, got)
			}
		})
	}
}
