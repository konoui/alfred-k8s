package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/konoui/alfred-k8s/pkg/kubectl"
	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

func SetupCmd(t *testing.T, cmd *cobra.Command, args []string) (outBuf, errBuf *bytes.Buffer) {
	t.Helper()

	// set global variables on behalf init()
	kubectl.TestDataBaseDir = "../pkg/kubectl/"
	k = kubectl.SetupKubectl(t, nil)
	awf = alfred.NewWorkflow()
	awf.EmptyWarning("There are no resources", "No matching")

	outBuf, errBuf = new(bytes.Buffer), new(bytes.Buffer)
	outStream, errStream = outBuf, errBuf
	cmd.SetOut(outStream)
	cmd.SetErr(errStream)
	awf.SetStreams(outStream, errStream)

	cmd.SetArgs(args)
	os.Args = append([]string{"dummy"}, args...)
	return outBuf, errBuf
}

func TestExecute(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		update bool
	}{
		{
			name: "list-available-commands",
			args: []string{
				"",
			},
		},
		{
			name: "empty-warning",
			args: []string{
				"xxxxxxxxxx",
			},
		},
		{
			name: "list-available-commands-for-fuzzy",
			args: []string{
				"no",
			},
		},
		{
			name: "list-nodes",
			args: []string{
				"node",
			},
		},
		{
			name: "list-pods",
			args: []string{
				"pod",
			},
		},
		{
			name: "list-services",
			args: []string{
				"svc",
			},
		},
		{
			name: "list-ingresses",
			args: []string{
				"ingress",
			},
		},
		{
			name: "list-deployments",
			args: []string{
				"deploy",
			},
		},
		{
			name: "list-contexts",
			args: []string{
				"context",
			},
		},
		{
			name: "list-namespaces",
			args: []string{
				"ns",
			},
		},
		{
			name: "list-crds",
			args: []string{
				"crd",
			},
		},
		{
			name: "list-base-pods",
			args: []string{
				"obj",
				"pod",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fmt.Sprintf("testdata/%s.json", tt.name)
			rootCmd := NewDefaultCmd()
			outBuf, _ := SetupCmd(t, rootCmd, tt.args)
			Execute(rootCmd)
			outGotData := outBuf.Bytes()

			if tt.update {
				if err := ioutil.WriteFile(f, outGotData, 0644); err != nil {
					t.Fatal(err)
				}
			}
			outWantData, err := ioutil.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}

			if diff := alfred.DiffScriptFilter(outWantData, outGotData); diff != "" {
				t.Errorf("-want +got\n%+v", diff)
			}
		})
	}
}
