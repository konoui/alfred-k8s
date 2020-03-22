package executor

// Executor is interface of command execution
type Executor interface {
	Exec(args ...string) (*Response, error)
}

// Response is command execution response
type Response struct {
	ExitCode int
	Stdout   []byte
	Stderr   []byte
}
