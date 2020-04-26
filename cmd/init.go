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
	k         *kubectl.Kubectl
	awf       *alfred.Workflow
	cacheTime time.Duration
	cacheDir  = os.TempDir()
)

const (
	cacheSuffix              = "-alfred-k8s.cache"
	cacheNamespacedPrefix    = "namespaced-"
	cacheNonNamespacedPrefix = "non-namespaced-"
)

// decide next action for workflow filter
const (
	nextActionKey   = "nextAction"
	nextActionCmd   = "cmd"
	nextActionShell = "shell"
)

func init() {
	awf = alfred.NewWorkflow()
	awf.SetOut(outStream)
	awf.EmptyWarning("There are no resources", "No matching")
	awf.SetCacheSuffix(cacheSuffix)
	err := awf.SetCacheDir(cacheDir)
	exitWith(err)

	c, err := newConfig()
	exitWith(err)

	var binOpt, pluginPathOpt kubectl.Option
	if c.Kubectl.Bin != "" {
		binOpt = kubectl.OptionBinary(c.Kubectl.Bin)
	}
	if paths := c.Kubectl.PluginPaths; len(paths) > 0 {
		path := strings.Join(paths, ":")
		pluginPathOpt = kubectl.OptionPluginPath(path)
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

	k, err = kubectl.New(binOpt, pluginPathOpt)
	exitWith(err)
}
