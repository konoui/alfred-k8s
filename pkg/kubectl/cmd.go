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

var timeout = 10 * time.Second

// CmdResponse is command execution response
type CmdResponse struct {
	exitCode int
	err      error
	stdout   []byte
	stderr   []byte
}

func generateKubectlCmd(arg string) string {
	return fmt.Sprintf("%s %s", findKubeclt(), arg)
}

func findKubeclt() string {
	// FIXME
	return "/usr/local/bin/kubectl"
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

// Execute run command
func Execute(command string) *CmdResponse {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmds, err := shellwords.Parse(command)
	if err != nil {
		return &CmdResponse{}
	}

	var stdout, stderr bytes.Buffer
	// FIXME check length of cmds[1:]
	cmd := exec.CommandContext(ctx, cmds[0], cmds[1:]...)
	cmd.Env = os.Environ()
	// FIXME
	cmd.Env = append(cmd.Env, "PATH=/usr/local/bin/")
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
