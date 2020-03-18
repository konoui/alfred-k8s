package kubectl

import (
	"fmt"
)

// GetCurrentContext return current configuration
func GetCurrentContext() (string, error) {
	cmd := generateKubectlCmd("config current-context")
	resp := Execute(cmd)
	return <-resp.Readline(), resp.err
}

// SetNamespace configure namepsace
func SetNamespace(context, ns string) error {
	arg := fmt.Sprintf("config set-context %s --namespace=%s", context, ns)
	cmd := generateKubectlCmd(arg)
	resp := Execute(cmd)
	return resp.err
}
