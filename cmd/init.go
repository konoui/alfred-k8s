package cmd

import (
	"os"
	"strings"
	"time"

	"github.com/konoui/alfred-k8s/cmd/config"
	"github.com/konoui/alfred-k8s/cmd/portforwardcmd"
	"github.com/konoui/alfred-k8s/cmd/utils"
	"github.com/konoui/alfred-k8s/pkg/kubectl"
	"github.com/konoui/go-alfred"
)

var (
	k            *kubectl.Kubectl
	awf          *alfred.Workflow
	cacheTime    time.Duration
	cacheDir     = os.TempDir()
	experimental = false
)

const (
	cacheSuffix   = "-alfred-k8s.cache"
	emptyTitle    = "There are no resources"
	emptySubTitle = "No matching"
)

func init() {
	awf = alfred.NewWorkflow()
	awf.SetOut(outStream)
	awf.SetErr(errStream)
	awf.EmptyWarning(emptyTitle, emptySubTitle)
	awf.SetCacheSuffix(cacheSuffix)
	err := awf.SetCacheDir(cacheDir)
	exitWith(err)

	cfg, err := config.New()
	exitWith(err)

	initAlfredMod()

	var opts []kubectl.Option
	if cfg.Kubectl.Bin != "" {
		opts = append(opts, kubectl.OptionBinary(cfg.Kubectl.Bin))
	}
	if paths := cfg.Kubectl.PluginPaths; len(paths) > 0 {
		path := strings.Join(paths, ":")
		opts = append(opts, kubectl.OptionPluginPath(path))
	}
	cacheTime = cfg.TTL()

	k, err = kubectl.New(opts...)
	exitWith(err)
}

func initAlfredMod() {
	// FIXME exec Portforward is experimental
	if experimental {
		utils.GetPortForwardMod = portforwardcmd.GetExecPortForwardMod
		return
	}
	utils.GetPortForwardMod = portforwardcmd.GetCopyPortForwardMod
}

func exitWith(err error) {
	if err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	awf.Fatal("Fatal error occurs", err.Error())
}
