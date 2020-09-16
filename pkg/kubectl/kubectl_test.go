package kubectl

import (
	"os"
	"reflect"
	"strings"
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
				cmd: newCommand("/usr/local/bin/kubectl"),
				env: setPathEnv("/usr/local/bin/"),
			},
		},
		{
			name: "options value",
			options: []Option{
				OptionBinary(knownBinary),
				OptionPluginPaths([]string{knownBinPath}),
			},
			want: &Kubectl{
				cmd: newCommand(knownBinary),
				env: setPathEnv(knownBinPath),
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

func TestOptionPluginPaths(t *testing.T) {
	path := "/usr/local/test"
	key := "TEST_USER"
	value := "test"
	input := "$" + key + path
	want := value + path
	t.Run("expand env test", func(t *testing.T) {
		// Note before set ENV
		if err := os.Setenv(key, value); err != nil {
			t.Fatal(err)
		}

		k, err := New(OptionPluginPaths([]string{input}))
		if err != nil {
			t.Fatal(err)
		}

		var got string
		for _, e := range k.env {
			key := strings.SplitN(e, "=", 2)[0]
			value := strings.SplitN(e, "=", 2)[1]
			if key == "PATH" {
				first := strings.SplitN(value, ":", 2)[0]
				got = first
				break
			}
		}

		if got != want {
			t.Errorf("unexpected want: %v\ngot: %v", want, got)
		}
	})
}

func TestOptionBinary(t *testing.T) {
	path := "/ls"
	key := "TEST_BIN"
	value := "/bin"
	input := "$" + key + path
	want := newCommand(value + path)
	t.Run("expand env test", func(t *testing.T) {
		if err := os.Setenv(key, value); err != nil {
			t.Fatal(err)
		}

		k, err := New(OptionBinary(input))
		if err != nil {
			t.Fatal(err)
		}

		got := k.cmd
		if !reflect.DeepEqual(want, got) {
			t.Errorf("unexpected want: %v\ngot: %v", want, got)
		}
	})
}
