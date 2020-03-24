package kubectl

import (
	"github.com/konoui/alfred-k8s/pkg/executor"
)

// FakeExecutor is mock for test
type FakeExecutor struct {
	Impl FakeFunc
}

// FakeFunc is real implementation for mock test
// please contain stdout/stderr message in executor.Response.Stdout/Stderr
type FakeFunc func(args ...string) (*executor.Response, error)

// Exec is mock for test
// Exec return executor.Response and error in FakeFunc
func (e *FakeExecutor) Exec(args ...string) (*executor.Response, error) {
	return e.Impl(args...)
}

// NewFakeExecutor is mock for test
func NewFakeExecutor(impl FakeFunc) executor.Executor {
	return &FakeExecutor{
		Impl: impl,
	}
}

// OptionExecutor is mock for configuration of kubectl execution function for test
// If this option is set, OptionBinary must not set.
func OptionExecutor(e executor.Executor) Option {
	return func(k *Kubectl) error {
		k.cmd = e
		return nil
	}
}
