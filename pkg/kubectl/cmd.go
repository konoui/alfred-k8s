package kubectl

import (
	"bufio"
	"bytes"
	"context"
	"io"
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

func (c *Command) String() string {
	return c.bin
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

	resp, err := k.cmd.Exec(args, k.env)
	return &Response{
		exitCode: resp.ExitCode,
		stdout:   resp.Stdout,
		stderr:   resp.Stderr,
	}, errors.Wrapf(err, resp.Stderr.String())
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
func (c *Command) Exec(args, env []string) (*executor.Response, error) {
	var stdout, stderr bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, c.bin, args...) //nolint:gosec //nolint:gocritic
	cmd.Env = env
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
