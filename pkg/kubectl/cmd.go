package kubectl

import (
	"bufio"
	"bytes"
	"context"
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
func New(opts ...Option) *Kubectl {
	k := &Kubectl{
		bin:        "/usr/local/bin/kubectl",
		pluginPath: "/usr/local/bin/",
	}

	for _, opt := range opts {
		if err := opt(k); err != nil {
			panic(err)
		}
	}

	return k
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
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "PATH="+k.pluginPath)
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
	s := bytes.NewReader(c.stdout)
	scanner := bufio.NewScanner(s)
	var outStream = make(chan string)

	go func() {
		defer close(outStream)
		for scanner.Scan() {
			outStream <- scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			outStream <- err.Error()
		}
	}()

	return outStream
}
