package namespacecmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/konoui/alfred-k8s/cmd/rootcmd"
	"github.com/konoui/alfred-k8s/cmd/utils"
	"github.com/konoui/go-alfred"
	"github.com/peterbourgon/ff/v3/ffcli"
)

type Config struct {
	use        bool
	fs         *flag.FlagSet
	rootConfig *rootcmd.Config
}

const CmdName = "ns"

// New create a new cmd for namespace resource
func New(rootConfig *rootcmd.Config) *ffcli.Command {
	fs := flag.NewFlagSet(CmdName, flag.ContinueOnError)
	cfg := &Config{
		rootConfig: rootConfig,
		fs:         fs,
	}
	cfg.registerFlags()

	cmd := &ffcli.Command{
		Name:      CmdName,
		ShortHelp: "list namespaces",
		FlagSet:   fs,
		Exec: func(ctx context.Context, args []string) error {
			if cfg.use {
				return cfg.useNamespace()
			}
			return cfg.collectNamespaces()
		},
	}

	return cmd
}

func (cfg *Config) registerFlags() {
	cfg.fs.BoolVar(&cfg.use, utils.UseFlag, false, "use it")
}

func (cfg *Config) useNamespace() (err error) {
	ns := cfg.fs.Arg(0)
	if err = cfg.rootConfig.Kubeclt().UseNamespace(ns); err != nil {
		fmt.Fprintf(cfg.rootConfig.Stdout(), "Failed due to %s\n", err)
		return nil
	}
	fmt.Fprintf(cfg.rootConfig.Stdout(), "Success!!\n")
	return
}

func (cfg *Config) collectNamespaces() (err error) {
	namespaces, err := cfg.rootConfig.Kubeclt().GetNamespaces()
	if err != nil {
		return
	}
	for _, ns := range namespaces {
		title := ns.Name
		if ns.Current {
			title = fmt.Sprintf("[*] %s", ns.Name)
		}
		cfg.rootConfig.Awf().Append(
			alfred.NewItem().
				SetTitle(title).
				SetSubtitle(
					fmt.Sprintf("status [%s] age [%s]", ns.Status, ns.Age),
				).
				SetArg(ns.Name).
				SetMod(alfred.ModCtrl, utils.GetUseMod("ns", ns)),
		)
	}

	cfg.rootConfig.Awf().Filter(cfg.fs.Arg(0)).Output()
	return
}
