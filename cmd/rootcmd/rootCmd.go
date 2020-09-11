package rootcmd

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"time"

	"github.com/konoui/alfred-k8s/pkg/kubectl"
	"github.com/konoui/go-alfred"
	"github.com/peterbourgon/ff/v3/ffcli"
)

const (
	cacheSuffix   = "-alfred-k8s.cache"
	emptyTitle    = "There are no resources"
	emptySubTitle = "No matching"
)

type Config struct {
	outStream io.Writer
	errStream io.Writer
	k         *kubectl.Kubectl
	awf       *alfred.Workflow
	cacheTTL  time.Duration
	cacheDir  string
	KeyMaps   *KeyMaps
}

func NewConfig(out, errOut io.Writer) *Config {
	cacheDir := os.TempDir()
	awf := alfred.NewWorkflow()
	awf.SetOut(out)
	awf.SetErr(errOut)
	awf.EmptyWarning(emptyTitle, emptySubTitle)
	awf.SetCacheSuffix(cacheSuffix)
	err := awf.SetCacheDir(cacheDir)
	if err != nil {
		awf.Fatal("Fatal error occurs on initialization", err.Error())
	}

	cfgFile, err := newConfigFile()
	if err != nil {
		awf.Fatal("Fatal error occurs on initialization", err.Error())
	}

	var opts []kubectl.Option
	if cfgFile.Kubectl.Bin != "" {
		opts = append(opts, kubectl.OptionBinary(cfgFile.Kubectl.Bin))
	}
	if paths := cfgFile.Kubectl.PluginPaths; len(paths) > 0 {
		path := strings.Join(paths, ":")
		opts = append(opts, kubectl.OptionPluginPath(path))
	}

	k, err := kubectl.New(opts...)
	if err != nil {
		awf.Fatal("Fatal error occurs on initialization", err.Error())
	}

	return &Config{
		outStream: out,
		errStream: errOut,
		k:         k,
		awf:       awf,
		cacheDir:  cacheDir,
		cacheTTL:  cfgFile.cacheTTL(),
		KeyMaps:   &cfgFile.KeyMaps,
	}
}

func (cfg *Config) Stdout() io.Writer {
	return cfg.outStream
}

func (cfg *Config) Stderr() io.Writer {
	return cfg.errStream
}

func (cfg *Config) Kubeclt() *kubectl.Kubectl {
	return cfg.k
}

func (cfg *Config) SetKubeCtl(k *kubectl.Kubectl) {
	cfg.k = k
}

func (cfg *Config) Awf() *alfred.Workflow {
	return cfg.awf
}

func (cfg *Config) CacheTTL() time.Duration {
	return cfg.cacheTTL
}

func (cfg *Config) SetCacheTTL(t time.Duration) {
	cfg.cacheTTL = t
}

// New create a new cmd for root
func New() *ffcli.Command {
	rootCmd := &ffcli.Command{
		Name: "list available commands",
		Exec: func(ctx context.Context, args []string) error {
			return errors.New("not implemented")
		},
	}

	return rootCmd
}
