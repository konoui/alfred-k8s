package kubectl

import (
	"github.com/konoui/alfred-k8s/pkg/executor"
)

// FakeFunc is real implementation for mock test
type FakeFunc func(args ...string) (*executor.Response, error)

// FakeExecutor is mock for test
type FakeExecutor struct {
	Impl FakeFunc
}

// Exec is mock for test
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
