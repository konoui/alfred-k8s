package kubectl

import (
	"bufio"
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// Job is represent async job status
type Job struct {
	outCh chan string
	Pid   int `json:"Pid"`
}

// StartJob is async command execution
func (k *Kubectl) StartJob(ctx context.Context, args ...string) (resp *Job, errCh <-chan error) {
	kbin := k.cmd.(*Command)
	appendPathEnv(k.pluginPath)
	resp, errCh = startJob(ctx, kbin.bin, args...)
	return
}

// StartJob is async command execution
func startJob(ctx context.Context, name string, args ...string) (*Job, <-chan error) { //nolint:gocritic
	cmd := exec.Command(name, args...) //nolint:gosec //nolint:gocritic
	cmd.Env = os.Environ()
	errChan := make(chan error)
	job := &Job{
		Pid:   -1,
		outCh: make(chan string),
	}

	ready := make(chan struct{})
	go func() {
		defer close(errChan)
		defer close(ready)
		stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()
		if err := cmd.Start(); err != nil {
			errChan <- err
			return
		}

		// notify ready to set pid
		job.Pid = cmd.Process.Pid
		ready <- struct{}{}

		done := make(chan struct{})
		scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
		go func() {
			defer close(done)
			defer close(job.outCh)
			for scanner.Scan() {
				job.outCh <- scanner.Text()
			}
		}()

		select {
		case <-ctx.Done():
			_ = cmd.Process.Kill()
			errChan <- errors.New("job is canceled")
		case <-done:
			errChan <- cmd.Wait()
		}
	}()

	// wait for set pid or close chan
	<-ready
	return job, errChan
}

// ReadLine return stdout and stderr chan
func (j *Job) ReadLine() <-chan string {
	return j.outCh
}

// Terminate kill the job
func (j *Job) Terminate() error {
	p, err := os.FindProcess(j.Pid)
	if err != nil {
		return errors.Wrapf(err, "failed to find process %d", j.Pid)
	}
	return p.Kill()
}
