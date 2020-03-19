package cmd

import (
	"os"

	"github.com/konoui/alfred-k8s/pkg/kubectl"
	"github.com/konoui/go-alfred"
)

var k *kubectl.Kubectl
var awf *alfred.Workflow

// decide next action for workflow filter
const (
	nextActionKey = "nextAction"
	nextActionCmd = "cmd"
)

func init() {
	awf = alfred.NewWorkflow()
	// alfred script filter read from only stdout
	awf.SetStreams(outStream, outStream)
	awf.EmptyWarning("There are no resources", "No matching")

	c, err := newConfig()
	if err != nil {
		awf.Fatal("fatal error occurs", err.Error())
		awf.Output()
		os.Exit(255)
	}

	binOpt, pluginPathOpt := kubectl.OptionNone(), kubectl.OptionNone()
	if c.kubectl.bin != "" {
		binOpt = kubectl.OptionBinary(c.kubectl.bin)
	}
	if c.kubectl.pluginPath != "" {
		pluginPathOpt = kubectl.OptionPluginPath(c.kubectl.pluginPath)
	}

	k = kubectl.New(binOpt, pluginPathOpt)
}
