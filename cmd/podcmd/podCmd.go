package podcmd

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
	all        bool
	delete     bool
	namespace  string
	fs         *flag.FlagSet
	rootConfig *rootcmd.Config
}

const CmdName = "pod"

// New create a new cmd for pod resource
func New(rootConfig *rootcmd.Config) *ffcli.Command {
	fs := flag.NewFlagSet(CmdName, flag.ContinueOnError)
	cfg := &Config{
		rootConfig: rootConfig,
		fs:         fs,
	}
	cfg.registerFlags()

	cmd := &ffcli.Command{
		Name:      CmdName,
		ShortHelp: "list pods",
		FlagSet:   fs,
		Exec: func(ctx context.Context, args []string) error {
			if cfg.delete {
				return cfg.rootConfig.DeleteOutput(
					cfg,
					cfg.GetQuery(),
				)
			}
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
	cfg.fs.BoolVar(&cfg.delete, utils.DeleteFlag, false, "delete it")
	cfg.fs.BoolVar(&cfg.all, utils.AllNamespacesFlag, false, "in all namespaces")

}

func (cfg *Config) Collect() error {
	pods, err := cfg.rootConfig.Kubeclt().GetPods(cfg.all)
	if err != nil {
		return err
	}
	for _, p := range pods {
		title := utils.GetNamespacedResourceTitle(p)
		modMap := map[rootcmd.KeyMapKey]*alfred.Mod{
			rootcmd.CopyResourceKey: utils.GetCopyMod(
				fmt.Sprintf("ready [%s] status [%s] restarts [%s] ", p.Ready, p.Status, p.Restarts),
				p.Name,
			),
			rootcmd.DeleteResourceKey:  utils.GetDeleteMod(CmdName, p),
			rootcmd.CopySternKey:       utils.GetSternMod(p),
			rootcmd.CopyPortForwardKey: portforwardcmd.GetCopyPortForwardMod(cfg.rootConfig.Kubeclt(), CmdName, p),
			rootcmd.ExecPortForwardKey: portforwardcmd.GetExecPortForwardMod(cfg.rootConfig.Kubeclt(), CmdName, p),
		}
		enterMod, mods := rootcmd.MakeMods(&cfg.rootConfig.KeyMaps.PodKeyMap, modMap)
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

	return nil
}

func (cfg *Config) Delete() error {
	pod := cfg.fs.Arg(0)
	arg := fmt.Sprintf("delete %s %s", CmdName, pod)
	if cfg.namespace != "" {
		arg = fmt.Sprintf("%s --namespace %s", arg, cfg.namespace)
	}
	_, _ = cfg.rootConfig.Kubeclt().Execute(arg)
	return nil
}

func (cfg *Config) GetQuery() string {
	return cfg.fs.Arg(0)
}
