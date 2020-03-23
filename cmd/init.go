package cmd

import (
	"os"
	"strings"

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
	checkWithExit(err)

	binOpt, pluginPathOpt := kubectl.OptionNone(), kubectl.OptionNone()
	if c.Kubectl.Bin != "" {
		binOpt = kubectl.OptionBinary(c.Kubectl.Bin)
	}
	if paths := c.Kubectl.PluginPaths; len(paths) > 0 {
		path := strings.Join(paths, ":")
		pluginPathOpt = kubectl.OptionPluginPath(path)
	}

	k, err = kubectl.New(binOpt, pluginPathOpt)
	if err != nil {
		checkWithExit(err)
	}
}

func checkWithExit(err error) {
	if err != nil {
		awf.Fatal("fatal error occurs", err.Error())
		awf.Output()
		os.Exit(255)
	}
}
