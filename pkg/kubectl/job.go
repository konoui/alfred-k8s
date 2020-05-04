package kubectl

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

// Status is represent async job status
type Status struct {
	pidChan chan int
	outChan chan string
	pid     int
}

// StartJob is async command execution
func (k *Kubectl) StartJob(args ...string) (resp *Status, errChan <-chan error) {
	kbin := k.cmd.(*Command)
	appendPathEnv(k.pluginPath)
	resp, errChan = startJob(kbin.bin, args...)
	return
}

// StartJob is async command execution
func startJob(name string, args ...string) (*Status, <-chan error) { //nolint:gocritic
	cmd := exec.Command(name, args...) //nolint:gosec //nolint:gocritic
	cmd.Env = os.Environ()
	errChan := make(chan error)
	status := &Status{
		pid:     -1,
		pidChan: make(chan int, 1),
		outChan: make(chan string),
	}

	go func() {
		stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()
		defer close(errChan)
		if err := cmd.Start(); err != nil {
			status.pidChan <- 0
			errChan <- err
			return
		}

		status.pidChan <- cmd.Process.Pid
		close(status.pidChan)
		scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(status.outChan)
			for scanner.Scan() {
				status.outChan <- scanner.Text()
			}
		}()

		wg.Wait()
		errChan <- cmd.Wait()
	}()

	return status, errChan
}

// ReadLine return stdout and stderr chan
func (s *Status) ReadLine() <-chan string {
	return s.outChan
}

// TerminateJob kill the job
func (s *Status) TerminateJob() error {
	pid := s.Pid()
	return syscall.Kill(-pid, syscall.SIGTERM)
}

// Pid return pid
func (s *Status) Pid() int {
	pid, ok := <-s.pidChan
	// 1st call, set pid
	if ok {
		s.pid = pid
		return pid
	}
	// for 2nd call
	return s.pid
}
