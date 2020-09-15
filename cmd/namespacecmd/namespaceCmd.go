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
				return cfg.rootConfig.UseOutput(
					cfg,
					cfg.GetQuery(),
				)
			}
			return cfg.rootConfig.CollectOutput(
				cfg,
				cfg.GetQuery(),
				utils.GetCacheKey(CmdName, false),
			)
		},
	}

	return cmd
}

func (cfg *Config) registerFlags() {
	cfg.fs.BoolVar(&cfg.use, utils.UseFlag, false, "use it")
}

func (cfg *Config) Use(ns string) (err error) {
	_ = cfg.rootConfig.Kubeclt().UseNamespace(ns)
	return nil
}

func (cfg *Config) Collect() (err error) {
	namespaces, err := cfg.rootConfig.Kubeclt().GetNamespaces()
	if err != nil {
		return
	}
	for _, ns := range namespaces {
		title := ns.Name
		if ns.Current {
			title = fmt.Sprintf("[*] %s", ns.Name)
		}

		modMap := map[rootcmd.KeyMapKey]*alfred.Mod{
			rootcmd.CopyResourceKey: utils.GetCopyMod(
				fmt.Sprintf("status [%s] age [%s]", ns.Status, ns.Age),
				ns.Name,
			),
			rootcmd.UseResourceKey: utils.GetUseMod("ns", ns),
		}

		enterMod, mods := rootcmd.MakeMods(&cfg.rootConfig.KeyMaps.NamespaceKeyMap, modMap)
		arg := enterMod.Arg
		subtitle := enterMod.Subtitle
		vals := enterMod.Variables
		cfg.rootConfig.Awf().Append(
			alfred.NewItem().
				SetTitle(title).
				SetSubtitle(subtitle).
				SetArg(arg).
				SetVariables(vals).
				SetMods(mods),
		)
	}

	return
}

func (cfg *Config) GetQuery() string {
	return cfg.fs.Arg(0)
}
