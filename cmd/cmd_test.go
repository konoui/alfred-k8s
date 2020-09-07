package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/konoui/alfred-k8s/pkg/kubectl"
	"github.com/konoui/go-alfred"
)

func ExecuteDefaultCmd(t *testing.T, args []string) (outBuf, errBuf *bytes.Buffer) {
	outBuf, errBuf = SetupCmd(t, args)
	rootCmd := NewDefaultCmd()
	Execute(rootCmd)
	return outBuf, errBuf
}

func SetupCmd(t *testing.T, args []string) (outBuf, errBuf *bytes.Buffer) {
	t.Helper()

	// set global variables on behalf init()
	kubectl.TestDataBaseDir = "../pkg/kubectl/"
	cacheTime = 0 * time.Second
	k = kubectl.SetupKubectl(t, nil)
	awf = alfred.NewWorkflow()
	awf.EmptyWarning(emptyTitle, emptySubTitle)

	outBuf, errBuf = new(bytes.Buffer), new(bytes.Buffer)
	outStream, errStream = outBuf, errBuf
	awf.SetOut(outStream)

	os.Args = append([]string{"dummy"}, args...)
	return outBuf, errBuf
}

func TestListExecution(t *testing.T) {
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
			name: "list-invalid-flag",
			args: []string{
				"node",
				"-z",
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
			name: "list-pods-in-all-ns",
			args: []string{
				"pod",
				"-a",
			},
		},
		{
			name: "list-deployments",
			args: []string{
				"deploy",
			},
		},
		{
			name: "list-deployments-in-all-ns",
			args: []string{
				"deploy",
				"-a",
			},
		},
		{
			name: "list-services",
			args: []string{
				"svc",
			},
		},
		{
			name: "list-services-in-all-ns",
			args: []string{
				"svc",
				"-a",
			},
		},
		{
			name: "list-ingresses",
			args: []string{
				"ingress",
			},
		},
		{
			name: "list-ingresses-in-all-ns",
			args: []string{
				"ingress",
				"-a",
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
			name: "list-base-pods",
			args: []string{
				"obj",
				"po",
			},
		},
		{
			name: "list-base-pods-in-all-ns",
			args: []string{
				"obj",
				"-a",
				"pod",
			},
		},
		{
			name: "list-base-pods-in-all-ns-with-fuzzy",
			args: []string{
				"obj",
				"-a",
				"pod",
				"DUMMY-ARG",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outBuf, _ := ExecuteDefaultCmd(t, tt.args)
			outGotData := outBuf.Bytes()

			f := fmt.Sprintf("testdata/%s.json", tt.name)
			if tt.update {
				if err := ioutil.WriteFile(f, outGotData, 0600); err != nil {
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

func TestUseDeleteExecution(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		update bool
	}{
		{
			name: "use-dummy-context",
			args: []string{
				"context",
				"--use",
				"dummy",
			},
		},
		{
			name: "delete-dummy-context",
			args: []string{
				"context",
				"--delete",
				"dummy",
			},
		},
		{
			name: "use-dummy-namespace",
			args: []string{
				"ns",
				"--use",
				"dummy",
			},
		},
		{
			name: "delete-dummy-pod",
			args: []string{
				"pod",
				"--delete",
				"dummy",
			},
		},
		{
			name: "delete-dummy-pod-in-all-ns",
			args: []string{
				"pod",
				"-a",
				"--delete",
				"dummy",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outBuf, _ := ExecuteDefaultCmd(t, tt.args)
			outGotData := outBuf.Bytes()

			f := fmt.Sprintf("testdata/%s.txt", tt.name)
			if tt.update {
				if err := ioutil.WriteFile(f, outGotData, 0600); err != nil {
					t.Fatal(err)
				}
			}
			outWantData, err := ioutil.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(outWantData, outGotData) {
				t.Errorf("want: %v\ngot: %v", string(outWantData), string(outGotData))
			}
		})
	}
}
