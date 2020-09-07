package nodecmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/konoui/alfred-k8s/cmd/rootcmd"
	"github.com/konoui/go-alfred"
	"github.com/peterbourgon/ff/v3/ffcli"
)

type Config struct {
	fs         *flag.FlagSet
	rootConfig *rootcmd.Config
}

const CmdName = "node"

// New create a new cmd for node resource
func New(rootConfig *rootcmd.Config) *ffcli.Command {
	fs := flag.NewFlagSet(CmdName, flag.ContinueOnError)
	cfg := &Config{
		rootConfig: rootConfig,
		fs:         fs,
	}

	cmd := &ffcli.Command{
		Name:      CmdName,
		ShortHelp: "list nodes",
		FlagSet:   fs,
		Exec: func(ctx context.Context, args []string) error {
			return cfg.collectNodes()
		},
	}
	return cmd
}

func (cfg *Config) collectNodes() (err error) {
	nodes, err := cfg.rootConfig.Kubeclt().GetNodes()
	if err != nil {
		return
	}
	for _, n := range nodes {
		cfg.rootConfig.Awf().Append(
			alfred.NewItem().
				SetTitle(n.Name).
				SetSubtitle(
					fmt.Sprintf("status [%s] version [%s]", n.Status, n.Version),
				).
				SetArg(n.Name),
		)
	}

	cfg.rootConfig.Awf().Filter(cfg.fs.Arg(0)).Output()
	return
}
