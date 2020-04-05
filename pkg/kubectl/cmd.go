package kubectl

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/konoui/alfred-k8s/pkg/executor"
	"github.com/mattn/go-shellwords"
	"github.com/pkg/errors"
)

// Command is binary path and env to execute
type Command struct {
	bin string
}

// Response is output of kubectl command
type Response struct {
	exitCode int
	stdout   []byte
	stderr   []byte
}

func newCommand(bin string) executor.Executor {
	return &Command{
		bin: bin,
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
	}, errors.Wrapf(err, string(resp.Stderr))
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
func (c *Response) Readline() <-chan string {
	return c.ReadlineContext(context.Background())
}

// ReadlineContext return stdout chan
func (c *Response) ReadlineContext(ctx context.Context) <-chan string {
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

// Exec is implementation of command execution
func (c *Command) Exec(args ...string) (*executor.Response, error) {
	var stdout, stderr bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
			Stdout:   stdout.Bytes(),
			Stderr:   stderr.Bytes(),
		}, err
	}

	return &executor.Response{
		ExitCode: 0,
		Stdout:   stdout.Bytes(),
		Stderr:   stderr.Bytes(),
	}, nil
}
