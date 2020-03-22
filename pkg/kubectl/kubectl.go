package kubectl

import (
	"os/exec"

	"github.com/konoui/alfred-k8s/pkg/executor"
)

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
		if _, err := exec.LookPath(bin); err != nil {
			return err
		}
		k.cmd = newCommand(bin)
		return nil
	}
}

// OptionPluginPath is configuration of kubectl plugin path.
// e.g.) authentication command path
func OptionPluginPath(path string) Option {
	return func(k *Kubectl) error {
		k.pluginPath = path
		return nil
	}
}

// OptionNone noop
func OptionNone() Option {
	return func(k *Kubectl) error {
		return nil
	}
}
