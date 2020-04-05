package kubectl

import (
	"os"
	"os/exec"

	"github.com/konoui/alfred-k8s/pkg/executor"
)

const allNamespaceFlag = "--all-namespaces"

// Kubectl is configuration of binary paths
type Kubectl struct {
	cmd        executor.Executor
	pluginPath string
}

// Option is the type to replace default parameters.
type Option func(k *Kubectl) error

// New create kubectl instance
func New(opts ...Option) (*Kubectl, error) {
	k := &Kubectl{
		cmd:        newCommand("/usr/local/bin/kubectl"),
		pluginPath: "/usr/local/bin/",
	}

	for _, opt := range opts {
		if err := opt(k); err != nil {
			return nil, err
		}
	}

	return k, nil
}

// OptionBinary is configuration of kubectl absolute path
func OptionBinary(bin string) Option {
	return func(k *Kubectl) error {
		kctl := os.ExpandEnv(bin)
		if _, err := exec.LookPath(kctl); err != nil {
			return err
		}
		k.cmd = newCommand(kctl)
		return nil
	}
}

// OptionPluginPath is configuration of kubectl plugin path.
// e.g.) authentication command path
func OptionPluginPath(path string) Option {
	return func(k *Kubectl) error {
		// Replace ${HOME} with abs path
		k.pluginPath = os.ExpandEnv(path)
		return nil
	}
}
