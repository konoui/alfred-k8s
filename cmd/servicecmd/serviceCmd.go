package servicecmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/konoui/alfred-k8s/cmd/portforwardcmd"
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

const CmdName = "svc"

// New create a new cmd for service resource
func New(rootConfig *rootcmd.Config) *ffcli.Command {
	fs := flag.NewFlagSet(CmdName, flag.ContinueOnError)
	cfg := &Config{
		fs:         fs,
		rootConfig: rootConfig,
	}
	cfg.registerFlags()

	cmd := &ffcli.Command{
		Name:      CmdName,
		ShortHelp: "list services",
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
	svcs, err := cfg.rootConfig.Kubeclt().GetServices(cfg.all)
	if err != nil {
		return
	}
	for _, s := range svcs {
		title := utils.GetNamespacedResourceTitle(s)

		modMap := map[rootcmd.KeyMapKey]*alfred.Mod{
			rootcmd.CopyResourceKey: utils.GetCopyMod(
				fmt.Sprintf("cluster-ip [%s] external-ip [%s] ports [%s]", s.ClusterIP, s.ExternalIP, s.Ports),
				s.Name,
			),
			rootcmd.CopySternKey:       utils.GetSternMod(s),
			rootcmd.CopyPortForwardKey: portforwardcmd.GetCopyPortForwardMod(cfg.rootConfig.Kubeclt(), CmdName, s),
			rootcmd.ExecPortForwardKey: portforwardcmd.GetExecPortForwardMod(cfg.rootConfig.Kubeclt(), CmdName, s),
		}
		enterMod, mods := rootcmd.MakeMods(&cfg.rootConfig.KeyMaps.ServiceKeyMap, modMap)
		subtitle := enterMod.Subtitle
		arg := enterMod.Arg
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
