package basecmd

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

const CmdName = "obj"

// New create a new cmd for resources not supported
func New(rootConfig *rootcmd.Config) *ffcli.Command {
	fs := flag.NewFlagSet(CmdName, flag.ContinueOnError)
	cfg := &Config{
		rootConfig: rootConfig,
		fs:         fs,
	}
	cfg.registerFlags()

	cmd := &ffcli.Command{
		Name:      CmdName,
		ShortHelp: "list specific resources",
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
	name := cfg.fs.Arg(0)
	reses, err := cfg.rootConfig.Kubeclt().GetBaseResources(name, cfg.all)
	if err != nil {
		return
	}

	for _, r := range reses {
		title := utils.GetNamespacedResourceTitle(r)
		cfg.rootConfig.Awf().Append(
			alfred.NewItem().
				SetTitle(title).
				SetSubtitle(fmt.Sprintf("age [%s]", r.Age)).
				SetArg(r.Name),
		)
	}

	return
}

func (cfg *Config) GetQuery() string {
	return cfg.fs.Arg(1)
}
