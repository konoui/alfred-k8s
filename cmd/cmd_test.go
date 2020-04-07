package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/konoui/alfred-k8s/pkg/kubectl"
	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

func SetupCmd(t *testing.T, cmd *cobra.Command, args []string) (outBuf, errBuf *bytes.Buffer) {
	t.Helper()

	// set kubectl instance to global variable `k`
	k = kubectl.SetupKubectl(t, nil)
	outBuf, errBuf = new(bytes.Buffer), new(bytes.Buffer)
	outStream, errStream = outBuf, errBuf
	cmd.SetOut(outStream)
	cmd.SetErr(errStream)
	// need to set streams as init() is not called for the tests
	awf.SetStreams(outStream, errStream)

	cmd.SetArgs(args)
	os.Args = append([]string{"dummy"}, args...)
	return outBuf, errBuf
}

func TestExecute(t *testing.T) {
	type want struct {
		filepath string
		errMsg   string
	}
	tests := []struct {
		name string
		args []string
		want want
		cmd  *cobra.Command
	}{
		{
			name: "show available commands when no input",
			want: want{
				filepath: "testdata/list-available-commands.json",
			},
			args: []string{
				"",
			},
			cmd: NewDefaultCmd(),
		},
		{
			name: "match no sub command when invalid sub command",
			want: want{
				filepath: "testdata/empty-warning.json",
			},
			args: []string{
				"xxxxxxxxxx",
			},
			cmd: NewDefaultCmd(),
		},
		{
			name: "match some sub commands for fuzzy search",
			want: want{
				filepath: "testdata/list-available-commands-for-fuzzy.json",
			},
			args: []string{
				"no",
			},
			cmd: NewDefaultCmd(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outWantData, err := ioutil.ReadFile(tt.want.filepath)
			if err != nil {
				t.Fatal(err)
			}

			outBuf, errBuf := SetupCmd(t, tt.cmd, tt.args)
			Execute(tt.cmd)
			outGotData := outBuf.Bytes()
			errGotData := errBuf.Bytes()

			if diff := alfred.DiffScriptFilter(outWantData, outGotData); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}

			errWant := tt.want.errMsg
			if errWant != string(errGotData) {
				t.Errorf("want: %+v\n, got: %+v\n", errWant, string(errGotData))
			}
		})
	}
}
