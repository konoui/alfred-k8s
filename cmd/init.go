package cmd

import (
	"os"
	"strings"
	"time"

	"github.com/konoui/alfred-k8s/pkg/kubectl"
	"github.com/konoui/go-alfred"
)

const defaulCacheValue = 70

var (
	k                 *kubectl.Kubectl
	awf               *alfred.Workflow
	cacheTime         time.Duration
	cacheDir          = os.TempDir()
	experimental      = false
	getPortForwardMod func(string, interface{}) *alfred.Mod
)

const (
	cacheSuffix   = "-alfred-k8s.cache"
	emptyTitle    = "There are no resources"
	emptySubTitle = "No matching"
)

// decide next action for workflow filter
const (
	nextActionKey   = "nextAction"
	nextActionCmd   = "cmd"
	nextActionShell = "shell"
	nextActionJob   = "job"
)

func init() {
	awf = alfred.NewWorkflow()
	awf.SetOut(outStream)
	awf.SetErr(errStream)
	awf.EmptyWarning(emptyTitle, emptySubTitle)
	awf.SetCacheSuffix(cacheSuffix)
	err := awf.SetCacheDir(cacheDir)
	exitWith(err)

	c, err := newConfig()
	exitWith(err)

	initAlfredMod()

	var opts []kubectl.Option
	if c.Kubectl.Bin != "" {
		opts = append(opts, kubectl.OptionBinary(c.Kubectl.Bin))
	}
	if paths := c.Kubectl.PluginPaths; len(paths) > 0 {
		path := strings.Join(paths, ":")
		opts = append(opts, kubectl.OptionPluginPath(path))
	}
	// if minus value, disable cache. if zero value, set default cache time
	maxAge := c.CacheTimeSecond
	switch {
	case maxAge == 0:
		cacheTime = defaulCacheValue * time.Second
	case maxAge < 0:
		cacheTime = 0 * time.Second
	default:
		cacheTime = time.Duration(maxAge) * time.Second
	}

	k, err = kubectl.New(opts...)
	exitWith(err)
}

func initAlfredMod() {
	// FIXME exec Portforward is experimental
	if experimental {
		getPortForwardMod = getExecPortForwardMod
		return
	}
	getPortForwardMod = getCopyPortForwardMod
}
