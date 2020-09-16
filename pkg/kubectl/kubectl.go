package kubectl

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/konoui/alfred-k8s/pkg/executor"
)

const (
	allNamespaceFlag = "--all-namespaces"
)

// Kubectl is configuration of binary paths
type Kubectl struct {
	cmd executor.Executor
	env []string
}

// Option is the type to replace default parameters.
type Option func(k *Kubectl) error

// New create kubectl instance
func New(opts ...Option) (*Kubectl, error) {
	k := &Kubectl{
		cmd: newCommand("/usr/local/bin/kubectl"),
		env: setPathEnv("/usr/local/bin/"),
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
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

// OptionPluginPaths is configuration of kubectl plugin path.
// e.g.) authentication command path
func OptionPluginPaths(paths []string) Option {
	return func(k *Kubectl) error {
		for _, path := range paths {
			// Replace ${HOME} with abs path
			setPathEnv(os.ExpandEnv(path))
		}
		k.env = os.Environ()

		return nil
	}
}

func setPathEnv(pluginPath string) []string {
	key := "PATH"
	path := os.Getenv(key)
	if path == "" {
		os.Setenv(key, pluginPath)
		return os.Environ()
	}
	if strings.Contains(path, pluginPath) {
		return os.Environ()
	}
	os.Setenv(key, fmt.Sprintf("%s:%s", pluginPath, path))
	return os.Environ()
}

// GetKubectlCommandEnv return commands and env
func (k *Kubectl) GetKubectlCommandEnv(args []string) (cmds, env []string) {
	bin := fmt.Sprintf("%s", k.cmd)
	cmds = append([]string{bin}, args...)
	env = k.env
	return
}
