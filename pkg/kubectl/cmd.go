package kubectl

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/mattn/go-shellwords"
)

// CmdResponse is command execution response
type CmdResponse struct {
	exitCode int
	err      error
	stdout   []byte
	stderr   []byte
}

// Kubectl configuration of kubectl command
type Kubectl struct {
	bin        string
	pluginPath string
}

// Option is the type to replace default parameters.
type Option func(k *Kubectl) error

// OptionBinary is configuration of kubectl absolute path
func OptionBinary(bin string) Option {
	return func(k *Kubectl) error {
		if _, err := exec.LookPath(bin); err != nil {
			return err
		}
		k.bin = bin
		return nil
	}
}

// OptionPluginPath is configuration of kubectl plugin path.
// e.g.) authentication command path
func OptionPluginPath(path string) Option {
	return func(k *Kubectl) error {
		k.pluginPath = path
		return nil
	}
}

// OptionNone noop
func OptionNone() Option {
	return func(k *Kubectl) error {
		return nil
	}
}

// New create kubectl instance
func New(opts ...Option) (*Kubectl, error) {
	k := &Kubectl{
		bin:        "/usr/local/bin/kubectl",
		pluginPath: "/usr/local/bin/",
	}

	for _, opt := range opts {
		if err := opt(k); err != nil {
			return nil, err
		}
	}

	return k, nil
}

func overridePathEnv(addPath string) {
	path := os.Getenv("PATH")
	if path == "" {
		os.Setenv("PATH", addPath)
		return
	}
	os.Setenv("PATH", fmt.Sprintf("%s:%s", addPath, path))
}

// Execute run command
func (k *Kubectl) Execute(arg string) *CmdResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	args, err := shellwords.Parse(arg)
	if err != nil {
		return &CmdResponse{}
	}

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, k.bin, args...)
	overridePathEnv(k.pluginPath)
	cmd.Env = os.Environ()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		exitCode := 255
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		}

		return &CmdResponse{
			exitCode: exitCode,
			err:      err,
			stdout:   stdout.Bytes(),
			stderr:   stderr.Bytes(),
		}
	}

	return &CmdResponse{
		exitCode: 0,
		err:      nil,
		stdout:   stdout.Bytes(),
		stderr:   stderr.Bytes(),
	}
}

// Readline return stdout chan
func (c *CmdResponse) Readline() <-chan string {
	return c.ReadlineContext(context.Background())
}

// ReadlineContext return stdout chan
func (c *CmdResponse) ReadlineContext(ctx context.Context) <-chan string {
	s := bytes.NewReader(c.stdout)
	scanner := bufio.NewScanner(s)
	var outStream = make(chan string)

	go func(ctx context.Context) {
		defer close(outStream)
		for scanner.Scan() {
			select {
			case outStream <- scanner.Text():
			case <-ctx.Done():
				return
			}
		}

		if err := scanner.Err(); err != nil {
			select {
			case outStream <- err.Error():
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	return outStream
}
