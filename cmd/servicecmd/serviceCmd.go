package servicecmd

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
			return cfg.collectServices()
		},
	}

	return cmd
}

func (cfg *Config) registerFlags() {
	cfg.fs.BoolVar(&cfg.all, utils.AllNamespacesFlag, false, "in all namespaces")
}

func (cfg *Config) collectServices() (err error) {
	svcs, err := cfg.rootConfig.Kubeclt().GetServices(cfg.all)
	if err != nil {
		return
	}
	for _, s := range svcs {
		title := utils.GetNamespacedResourceTitle(s)
		cfg.rootConfig.Awf().Append(
			alfred.NewItem().
				SetTitle(title).
				SetSubtitle(
					fmt.Sprintf("cluster-ip [%s] external-ip [%s] ports [%s]", s.ClusterIP, s.ExternalIP, s.Ports),
				).
				SetArg(s.Name).
				SetMod(alfred.ModShift, utils.GetSternMod(s)).
				SetMod(alfred.ModAlt, utils.GetPortForwardMod(cfg.rootConfig.Kubeclt(), CmdName, s)),
		)
	}

	cfg.rootConfig.Awf().Filter(cfg.fs.Arg(0)).Output()
	return
}
