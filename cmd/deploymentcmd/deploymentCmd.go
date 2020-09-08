package deploymentcmd

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
	fs         *flag.FlagSet
	rootConfig *rootcmd.Config
}

const CmdName = "deploy"

// New create a new cmd for deployment resource
func New(rootConfig *rootcmd.Config) *ffcli.Command {
	fs := flag.NewFlagSet(CmdName, flag.ContinueOnError)
	cfg := &Config{
		rootConfig: rootConfig,
		fs:         fs,
	}
	cfg.registerFlags()

	cmd := &ffcli.Command{
		Name:      CmdName,
		ShortHelp: "list deployments",
		FlagSet:   fs,
		Exec: func(ctx context.Context, args []string) error {
			return cfg.collectDeployments()
		},
	}

	return cmd
}

func (cfg *Config) registerFlags() {
	cfg.fs.BoolVar(&cfg.all, utils.AllNamespacesFlag, false, "in all namespaces")
}

func (cfg *Config) collectDeployments() (err error) {
	deps, err := cfg.rootConfig.Kubeclt().GetDeployments(cfg.all)
	if err != nil {
		return
	}
	for _, d := range deps {
		title := utils.GetNamespacedResourceTitle(d)

		modMap := map[rootcmd.KeyMapKey]*alfred.Mod{
			rootcmd.CopyResourceKey: utils.GetCopyMod(
				fmt.Sprintf("ready [%s] up-to-date [%s] available [%s]", d.Ready, d.UpToDate, d.Available),
				d.Name,
			),
			rootcmd.CopySternKey:       utils.GetSternMod(d),
			rootcmd.CopyPortForwardKey: portforwardcmd.GetCopyPortForwardMod(cfg.rootConfig.Kubeclt(), CmdName, d),
			rootcmd.ExecPortForwardKey: portforwardcmd.GetExecPortForwardMod(cfg.rootConfig.Kubeclt(), CmdName, d),
		}
		enterMod, mods := rootcmd.MakeMods(&cfg.rootConfig.KeyMaps.DeploymentKeyMap, modMap)
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

	cfg.rootConfig.Awf().Filter(cfg.fs.Arg(0)).Output()
	return
}
