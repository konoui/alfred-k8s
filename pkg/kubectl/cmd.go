package kubectl

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/konoui/alfred-k8s/pkg/executor"
	"github.com/mattn/go-shellwords"
	"github.com/pkg/errors"
)

// Command is binary path and env to execute
type Command struct {
	bin     string
	timeout time.Duration
}

// Response is output of kubectl command
type Response struct {
	exitCode int
	stdout   io.Reader
	stderr   io.Reader
}

func newCommand(bin string) executor.Executor {
	return &Command{
		bin:     bin,
		timeout: 10 * time.Second,
	}
}

// Execute kubectl command
func (k *Kubectl) Execute(arg string) (*Response, error) {
	args, err := shellwords.Parse(arg)
	if err != nil {
		return &Response{}, err
	}

	appendPathEnv(k.pluginPath)
	resp, err := k.cmd.Exec(args...)
	return &Response{
		exitCode: resp.ExitCode,
		stdout:   resp.Stdout,
		stderr:   resp.Stderr,
	}, errors.Wrapf(err, resp.Stderr.String())
}

func appendPathEnv(addPath string) {
	path := os.Getenv("PATH")
	if path == "" {
		os.Setenv("PATH", addPath)
		return
	}
	os.Setenv("PATH", fmt.Sprintf("%s:%s", addPath, path))
}

// Readline return stdout chan
func (r *Response) Readline() <-chan string {
	return r.ReadlineContext(context.Background())
}

// ReadlineContext return stdout chan
func (r *Response) ReadlineContext(ctx context.Context) <-chan string {
	if r.stdout == nil {
		r.stdout = new(bytes.Buffer)
	}

	outchan := make(chan string)
	scanner := bufio.NewScanner(r.stdout)
	go func(ctx context.Context) {
		defer close(outchan)
		for scanner.Scan() {
			select {
			case outchan <- scanner.Text():
			case <-ctx.Done():
				return
			}
		}

		if err := scanner.Err(); err != nil {
			select {
			// FIXME
			case outchan <- err.Error():
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	return outchan
}

// Exec is implementation of command execution
func (c *Command) Exec(args ...string) (*executor.Response, error) {
	var stdout, stderr bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, c.bin, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		exitCode := 255
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		}

		return &executor.Response{
			ExitCode: exitCode,
			Stdout:   &stdout,
			Stderr:   &stderr,
		}, err
	}

	return &executor.Response{
		ExitCode: 0,
		Stdout:   &stdout,
		Stderr:   &stderr,
	}, nil
}
