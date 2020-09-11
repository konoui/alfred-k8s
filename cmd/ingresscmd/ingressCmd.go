package ingresscmd

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
	fs         *flag.FlagSet
	all        bool
	rootConfig *rootcmd.Config
}

const CmdName = "ingress"

// New create a new cmd for ingress resource
func New(rootConfig *rootcmd.Config) *ffcli.Command {

	fs := flag.NewFlagSet(CmdName, flag.ContinueOnError)
	cfg := &Config{
		rootConfig: rootConfig,
		fs:         fs,
	}
	cfg.registerFlags()

	cmd := &ffcli.Command{
		Name:      CmdName,
		ShortHelp: "list ingresses",
		FlagSet:   fs,
		Exec: func(ctx context.Context, args []string) error {
			return cfg.rootConfig.CollectOutput(
				cfg,
				cfg.GetQuery(),
				utils.GetCacheKey(CmdName, cfg.all),
			)
		},
	}

	return cmd
}

func (cfg *Config) registerFlags() {
	cfg.fs.BoolVar(&cfg.all, utils.AllNamespacesFlag, false, "in all namespaces")
}

func (cfg *Config) Collect() (err error) {
	ingresses, err := cfg.rootConfig.Kubeclt().GetIngresses(cfg.all)
	if err != nil {
		return
	}
	for _, i := range ingresses {
		title := utils.GetNamespacedResourceTitle(i)
		cfg.rootConfig.Awf().Append(
			alfred.NewItem().
				SetTitle(title).
				SetSubtitle(
					fmt.Sprintf("host [%s] address [%s] ports [%s] ", i.Hosts, i.Address, i.Ports),
				).
				SetArg(i.Name).
				SetMod(alfred.ModCtrl,
					alfred.NewMod().
						SetSubtitle("copy ingress Address").
						SetArg(i.Address),
				),
		)
	}

	return
}

func (cfg *Config) GetQuery() string {
	return cfg.fs.Arg(0)
}
