package kubectl

import (
	"io/ioutil"
	"testing"

	"github.com/konoui/alfred-k8s/pkg/executor"
)

// FakeExecutor is mock for test
type FakeExecutor struct {
	Impl FakeFunc
	t    *testing.T
}

// FakeFunc is real implementation for mock test
// please contain stdout/stderr message in executor.Response.Stdout/Stderr
type FakeFunc func(t *testing.T, args ...string) (*executor.Response, error)

// Exec is mock for test
// Exec return executor.Response and error in FakeFunc
func (e *FakeExecutor) Exec(args ...string) (*executor.Response, error) {
	return e.Impl(e.t, args...)
}

// NewFakeExecutor is mock for test
func NewFakeExecutor(t *testing.T, impl FakeFunc) executor.Executor {
	return &FakeExecutor{
		Impl: impl,
		t:    t,
	}
}

// SetupKubectl is helper for mock test.
// Please pass test and function of FakeFunc type
func SetupKubectl(t *testing.T, fakeFunc FakeFunc) *Kubectl {
	t.Helper()
	e := NewFakeExecutor(t, fakeFunc)
	k, err := New(OptionExecutor(e))
	if err != nil {
		t.Fatal(err)
	}
	return k
}

// OptionExecutor is mock for configuration of kubectl execution function for test
// If this option is set, OptionBinary must not set.
func OptionExecutor(e executor.Executor) Option {
	return func(k *Kubectl) error {
		k.cmd = e
		return nil
	}
}

// GetByteFromTestFile get vlues as []byte from test file.
func GetByteFromTestFile(t *testing.T, path string) []byte {
	t.Helper()

	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to open file '%#v'", err)
	}
	return data
}

// GetStringFromTestFile get vlues as string from test file.
func GetStringFromTestFile(t *testing.T, path string) string {
	t.Helper()

	return string(GetByteFromTestFile(t, path))
}
