package rootcmd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewConfigFile(t *testing.T) {
	tests := []struct {
		name string
		want *configFile
	}{
		{
			name: "override default value.",
			want: &configFile{
				Kubectl: Kubectl{
					Bin:         "/bin/ls",
					PluginPaths: []string{"/bin", "/usr/bin"},
				},
				CacheTimeSecond: 2,
				KeyMaps: KeyMaps{
					ContextKeyMap: KeyMap{
						Enter: "copy",
						Ctrl:  "use",
						Cmd:   "delete",
					},
					NamespaceKeyMap: KeyMap{
						Enter: "copy",
						Ctrl:  "use",
					},
					PodKeyMap: KeyMap{
						Enter: "copy",
						Ctrl:  "delete",
						Shift: "stern",
					},
					DeploymentKeyMap: KeyMap{
						Enter: "copy",
					},
					ServiceKeyMap: KeyMap{
						Enter: "copy",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		c, err := newConfigFile()
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(c, tt.want) {
			t.Errorf("want: \n%+v, got: \n%+v", tt.want, c)
		}
	}
}
