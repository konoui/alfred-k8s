package kubectl

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/konoui/alfred-k8s/pkg/executor"
)

// TestDataBaseDir is directory path to kubectl/testdata
var TestDataBaseDir string

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
func SetupKubectl(t *testing.T, f FakeFunc) *Kubectl {
	t.Helper()
	if f == nil {
		f = FakeResourceFunc
	}
	e := NewFakeExecutor(t, f)
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

	if base := TestDataBaseDir; base != "" {
		path = filepath.Join(base, path)
	}

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

// GetCurrentContextFromtTestFile is helper to read test data
func GetCurrentContextFromtTestFile(t *testing.T) string {
	t.Helper()
	c := GetStringFromTestFile(t, "testdata/raw-current-context.txt")
	return strings.Replace(c, "\n", "", -1)
}

// GetCurrentNamespaceFromtTestFile is helper to read test data
func GetCurrentNamespaceFromtTestFile(t *testing.T) string {
	t.Helper()
	c := GetStringFromTestFile(t, "testdata/raw-current-namespace.txt")
	return strings.Replace(c, "\n", "", -1)
}

func hasQuery(args []string, i int, expect string) bool {
	return getQuery(args, i) == expect
}

func getQuery(args []string, i int) string {
	if len(args) > i {
		return args[i]
	}
	return ""
}

// FakeResourceFunc please call this for mock test
func FakeResourceFunc(t *testing.T, args ...string) (*executor.Response, error) {
	fs := []FakeFunc{
		FakePodFunc,
		FakeContextFunc,
		FakeNamespaceFunc,
		FakeDeploymentFunc,
		FakeServiceFunc,
		FakeNodeFunc,
		FakeIngressFunc,
		FakePodBaseResourceFunc,
		FakeCRDFunc,
	}
	for _, f := range fs {
		resp, err := f(t, args...)
		if err == nil {
			// found valid execution
			return resp, nil
		}
	}
	return &executor.Response{}, fmt.Errorf("match no execution function %s", args)
}

// FakeContextFunc behave kubectl get context
func FakeContextFunc(t *testing.T, args ...string) (*executor.Response, error) {
	rawDataCurrentContext := GetByteFromTestFile(t, "testdata/raw-current-context.txt")
	rawDataContexts := GetByteFromTestFile(t, "testdata/raw-contexts.txt")

	if hasQuery(args, 0, "config") {
		switch getQuery(args, 1) {
		case "current-context":
			return &executor.Response{
				Stdout: rawDataCurrentContext,
			}, nil
		case "view":
			return &executor.Response{
				Stdout: rawDataContexts,
			}, nil
		case "use-context", "set-context", "delete-context":
			return &executor.Response{}, nil
		}
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

// FakeCRDFunc behave kubectl get crd
func FakeCRDFunc(t *testing.T, args ...string) (*executor.Response, error) {
	rawDataCRDs := GetByteFromTestFile(t, "testdata/raw-crds.txt")

	if hasQuery(args, 1, "crd") {
		return &executor.Response{
			Stdout: rawDataCRDs,
		}, nil
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

// FakeDeploymentFunc behave kubectl get deploy
func FakeDeploymentFunc(t *testing.T, args ...string) (*executor.Response, error) {
	rawDataAllDeployments := GetByteFromTestFile(t, "testdata/raw-deployments-in-all-namespaces.txt")
	rawDataDeployments := GetByteFromTestFile(t, "testdata/raw-deployments.txt")

	if hasQuery(args, 0, "get") && hasQuery(args, 1, "deployment") {
		if hasQuery(args, 2, allNamespaceFlag) {
			return &executor.Response{
				Stdout: []byte(rawDataAllDeployments),
			}, nil
		}

		return &executor.Response{
			Stdout: []byte(rawDataDeployments),
		}, nil
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

// FakeIngressFunc behave kubectl get ingress
func FakeIngressFunc(t *testing.T, args ...string) (*executor.Response, error) {
	rawDataAllIngresses := GetByteFromTestFile(t, "testdata/raw-ingresses-in-all-namespaces.txt")
	rawDataIngresses := GetByteFromTestFile(t, "testdata/raw-ingresses.txt")

	if hasQuery(args, 0, "get") && hasQuery(args, 1, "ingress") {
		if hasQuery(args, 2, allNamespaceFlag) {
			return &executor.Response{
				Stdout: rawDataAllIngresses,
			}, nil
		}

		return &executor.Response{
			Stdout: rawDataIngresses,
		}, nil
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

// FakeNamespaceFunc behave kubectl get ns
func FakeNamespaceFunc(t *testing.T, args ...string) (*executor.Response, error) {
	rawDataNamespaces := GetByteFromTestFile(t, "testdata/raw-namespaces.txt")

	if hasQuery(args, 1, "namespace") {
		return &executor.Response{
			Stdout: rawDataNamespaces,
		}, nil
	}
	// Note: get current namespace and namespaces call context function
	return FakeContextFunc(t, args...)
}

// FakeNodeFunc behave kubectl get node
func FakeNodeFunc(t *testing.T, args ...string) (*executor.Response, error) {
	rawDataNodes := GetByteFromTestFile(t, "testdata/raw-nodes.txt")

	if hasQuery(args, 0, "get") && hasQuery(args, 1, "node") {
		return &executor.Response{
			Stdout: rawDataNodes,
		}, nil
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

// FakePodFunc behave kubectl get pod
func FakePodFunc(t *testing.T, args ...string) (*executor.Response, error) {
	rawDataAllPods := GetByteFromTestFile(t, "testdata/raw-pods-in-all-namespaces.txt")
	rawDataPods := GetByteFromTestFile(t, "testdata/raw-pods.txt")

	if hasQuery(args, 0, "get") && hasQuery(args, 1, "pod") {
		if hasQuery(args, 2, allNamespaceFlag) {
			return &executor.Response{
				Stdout: rawDataAllPods,
			}, nil
		}

		return &executor.Response{
			Stdout: rawDataPods,
		}, nil
	}
	if hasQuery(args, 0, "delete") && hasQuery(args, 1, "pod") {
		return &executor.Response{}, nil
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

// FakeServiceFunc behave kubectl get svc
func FakeServiceFunc(t *testing.T, args ...string) (*executor.Response, error) {
	rawDataAllServices := GetByteFromTestFile(t, "testdata/raw-services-in-all-namespaces.txt")
	rawDataServices := GetByteFromTestFile(t, "testdata/raw-services.txt")

	if hasQuery(args, 0, "get") && hasQuery(args, 1, "service") {
		if hasQuery(args, 2, allNamespaceFlag) {
			return &executor.Response{
				Stdout: rawDataAllServices,
			}, nil
		}

		return &executor.Response{
			Stdout: rawDataServices,
		}, nil
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

// FakePodBaseResourceFunc behave kubectl get po
func FakePodBaseResourceFunc(t *testing.T, args ...string) (*executor.Response, error) {
	rawDataAllPods := GetByteFromTestFile(t, "testdata/raw-pods-in-all-namespaces.txt")
	rawDataPods := GetByteFromTestFile(t, "testdata/raw-pods.txt")

	if hasQuery(args, 0, "get") && hasQuery(args, 1, "po") {
		if hasQuery(args, 2, allNamespaceFlag) {
			return &executor.Response{
				Stdout: []byte(rawDataAllPods),
			}, nil
		}

		return &executor.Response{
			Stdout: []byte(rawDataPods),
		}, nil
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}
