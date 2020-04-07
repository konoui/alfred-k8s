package kubectl

import (
	"fmt"
	"io/ioutil"
	"strings"
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

// GetCurrentContext is helper to read test data
func GetCurrentContext(t *testing.T) string {
	t.Helper()
	c := GetStringFromTestFile(t, "testdata/raw-current-context.txt")
	return strings.Replace(c, "\n", "", -1)
}

// GetCurrentNamespace is helper to read test data
func GetCurrentNamespace(t *testing.T) string {
	t.Helper()
	c := GetStringFromTestFile(t, "testdata/raw-current-namespace.txt")
	return strings.Replace(c, "\n", "", -1)
}

// FakePodBaseResourceFunc behave kubectl get pod
var FakePodBaseResourceFunc = func(t *testing.T, args ...string) (*executor.Response, error) {
	rawDataAllPods := GetByteFromTestFile(t, "testdata/raw-pods-in-all-namespaces.txt")
	rawDataPods := GetByteFromTestFile(t, "testdata/raw-pods.txt")
	if len(args) >= 1 {
		if args[0] == "get" && args[1] == "pod" {
			if len(args) == 3 && args[2] == allNamespaceFlag {
				return &executor.Response{
					Stdout: []byte(rawDataAllPods),
				}, nil
			}

			return &executor.Response{
				Stdout: []byte(rawDataPods),
			}, nil
		}
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

// FakeContextFunc behave kubectl get context
func FakeContextFunc(t *testing.T, args ...string) (*executor.Response, error) {
	rawDataCurrentContext := GetByteFromTestFile(t, "testdata/raw-current-context.txt")
	rawDataContexts := GetByteFromTestFile(t, "testdata/raw-contexts.txt")
	if len(args) >= 2 {
		if args[1] == "current-context" {
			return &executor.Response{
				Stdout: rawDataCurrentContext,
			}, nil
		}
		if args[1] == "view" {
			return &executor.Response{
				Stdout: rawDataContexts,
			}, nil
		}
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

// FakeCRDFunc behave kubectl get crd
func FakeCRDFunc(t *testing.T, args ...string) (*executor.Response, error) {
	rawDataCRDs := GetByteFromTestFile(t, "testdata/raw-crds.txt")
	if len(args) >= 1 {
		if args[1] == "crd" {
			return &executor.Response{
				Stdout: rawDataCRDs,
			}, nil
		}
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

// FakeDeploymentFunc behave kubectl get deploy
func FakeDeploymentFunc(t *testing.T, args ...string) (*executor.Response, error) {
	rawDataAllDeployments := GetByteFromTestFile(t, "testdata/raw-deployments-in-all-namespaces.txt")
	rawDataDeployments := GetByteFromTestFile(t, "testdata/raw-deployments.txt")
	if len(args) >= 4 {
		if args[1] == "deployment" && args[2] == allNamespaceFlag {
			return &executor.Response{
				Stdout: []byte(rawDataAllDeployments),
			}, nil
		}
	}
	if len(args) >= 2 {
		if args[1] == "deployment" {
			return &executor.Response{
				Stdout: []byte(rawDataDeployments),
			}, nil
		}
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

// FakeIngressFunc behave kubectl get ingress
func FakeIngressFunc(t *testing.T, args ...string) (*executor.Response, error) {
	rawDataAllIngresses := GetByteFromTestFile(t, "testdata/raw-ingresses-in-all-namespaces.txt")
	rawDataIngresses := GetByteFromTestFile(t, "testdata/raw-ingresses.txt")
	if len(args) >= 4 {
		if args[1] == "ingress" && args[2] == allNamespaceFlag {
			return &executor.Response{
				Stdout: rawDataAllIngresses,
			}, nil
		}
	}
	if len(args) >= 2 {
		if args[1] == "ingress" {
			return &executor.Response{
				Stdout: rawDataIngresses,
			}, nil
		}
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

// FakeNamespaceFunc behave kubectl get ns
func FakeNamespaceFunc(t *testing.T, args ...string) (*executor.Response, error) {
	rawDataNamespaces := GetByteFromTestFile(t, "testdata/raw-namespaces.txt")
	if len(args) >= 2 {
		if args[1] == "namespace" {
			return &executor.Response{
				Stdout: rawDataNamespaces,
			}, nil
		}
		// Note: get current namespace and namespaces call context function
		return FakeContextFunc(t, args...)
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

// FakeNodeFunc behave kubectl get node
func FakeNodeFunc(t *testing.T, args ...string) (*executor.Response, error) {
	rawDataNodes := GetByteFromTestFile(t, "testdata/raw-nodes.txt")
	if len(args) >= 2 {
		if args[1] == "node" {
			return &executor.Response{
				Stdout: rawDataNodes,
			}, nil
		}
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

// FakePodFunc behave kubectl get pod
func FakePodFunc(t *testing.T, args ...string) (*executor.Response, error) {
	rawDataAllPods := GetByteFromTestFile(t, "testdata/raw-pods-in-all-namespaces.txt")
	rawDataPods := GetByteFromTestFile(t, "testdata/raw-pods.txt")
	pod := "pod"
	if len(args) >= 3 {
		if args[1] == pod && args[2] == allNamespaceFlag {
			return &executor.Response{
				Stdout: rawDataAllPods,
			}, nil
		}
	}
	if len(args) >= 2 {
		if args[1] == pod {
			return &executor.Response{
				Stdout: rawDataPods,
			}, nil
		}
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

// FakeServiceFunc behave kubectl get svc
func FakeServiceFunc(t *testing.T, args ...string) (*executor.Response, error) {
	rawDataAllServices := GetByteFromTestFile(t, "testdata/raw-services-in-all-namespaces.txt")
	rawDataServices := GetByteFromTestFile(t, "testdata/raw-services.txt")
	if len(args) >= 4 {
		if args[1] == "service" && args[2] == allNamespaceFlag {
			return &executor.Response{
				Stdout: rawDataAllServices,
			}, nil
		}
	}
	if len(args) >= 2 {
		if args[1] == "service" {
			return &executor.Response{
				Stdout: rawDataServices,
			}, nil
		}
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}
