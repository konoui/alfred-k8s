package kubectl

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/goleak"
)

const (
	knownBinary  = "/bin/ls"
	knownBinPath = "/bin"
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
				bin:        "/usr/local/bin/kubectl",
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
				bin:        knownBinary,
				pluginPath: knownBinPath,
			},
			// TODO unexptected bin path case.
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := New(tt.options...); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestExecute(t *testing.T) {
	tests := []struct {
		name      string
		options   []Option
		cmdArg    string
		expectErr bool
	}{
		{
			name: "execute simple command",
			options: []Option{
				OptionBinary(knownBinary),
			},
			cmdArg: "-al",
		},
		// TODO output stderr command as expectErr case
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := New(tt.options...)
			if err != nil {
				t.Fatal(err)
			}

			resp := k.Execute(tt.cmdArg)
			if err = resp.err; err != nil {
				t.Error(err)
			}

			if resp.exitCode != 0 {
				t.Error("exit status is not zero")
			}

			if stderr := string(resp.stderr); stderr != "" {
				t.Errorf("stderr is not empty %s", stderr)
			}

			if string(resp.stdout) == "" {
				t.Error("stdout is empty")
			}
		})
	}
}

func TestReadLineContext(t *testing.T) {
	tests := []struct {
		name      string
		cmdResp   *CmdResponse
		expectErr bool
	}{
		{
			name: "multi lines",
			cmdResp: &CmdResponse{
				stdout:   []byte(fmt.Sprintln("stdout\nstdout\nstdout")),
				err:      nil,
				exitCode: 0,
			},
		},
		{
			name: "one line",
			cmdResp: &CmdResponse{
				stdout:   []byte(fmt.Sprintf("stdout")),
				err:      nil,
				exitCode: 0,
			},
		},
		{
			name:      "no stdout",
			expectErr: true,
			cmdResp:   &CmdResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			ctx, cancel := context.WithCancel(context.Background())
			got := <-tt.cmdResp.ReadlineContext(ctx)
			if !tt.expectErr && got == "" {
				t.Error("stdout is empty")
			}
			// avoid go routine leak
			cancel()
		})
	}
}

func TestReadLine(t *testing.T) {
	tests := []struct {
		name      string
		cmdResp   *CmdResponse
		expectErr bool
	}{
		{
			name: "one line",
			cmdResp: &CmdResponse{
				stdout:   []byte(fmt.Sprintf("stdout")),
				err:      nil,
				exitCode: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			got := <-tt.cmdResp.Readline()
			if !tt.expectErr && got == "" {
				t.Error("stdout is empty")
			}
		})
	}
}
