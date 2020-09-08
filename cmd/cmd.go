package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/konoui/alfred-k8s/cmd/basecmd"
	"github.com/konoui/alfred-k8s/cmd/contextcmd"
	"github.com/konoui/alfred-k8s/cmd/deploymentcmd"
	"github.com/konoui/alfred-k8s/cmd/ingresscmd"
	"github.com/konoui/alfred-k8s/cmd/namespacecmd"
	"github.com/konoui/alfred-k8s/cmd/nodecmd"
	"github.com/konoui/alfred-k8s/cmd/podcmd"
	"github.com/konoui/alfred-k8s/cmd/portforwardcmd"
	"github.com/konoui/alfred-k8s/cmd/rootcmd"
	"github.com/konoui/alfred-k8s/cmd/servicecmd"
	"github.com/konoui/alfred-k8s/cmd/utils"
	"github.com/konoui/alfred-k8s/cmd/versioncmd"
	"github.com/konoui/go-alfred"
	"github.com/peterbourgon/ff/v3/ffcli"
)

var (
	rootConfig *rootcmd.Config
)

func init() {
	rootConfig = rootcmd.NewConfig(os.Stdout, os.Stderr)
}

// Execute root cmd
func Execute(rootCmd *ffcli.Command) {
	if err := rootCmd.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		_ = collectAvailableSubCmds(
			subCmds(),
			os.Args[1:],
		)
	}
}

func subCmds() []*ffcli.Command {
	return []*ffcli.Command{
		versioncmd.New(rootConfig),
		contextcmd.New(rootConfig),
		podcmd.New(rootConfig),
		namespacecmd.New(rootConfig),
		nodecmd.New(rootConfig),
		deploymentcmd.New(rootConfig),
		servicecmd.New(rootConfig),
		ingresscmd.New(rootConfig),
		basecmd.New(rootConfig),
		portforwardcmd.New(rootConfig),
	}
}

// NewDefaultCmd create sub commands
func NewDefaultCmd() *ffcli.Command {
	rootCmd := rootcmd.New()
	rootCmd.Subcommands = subCmds()
	return rootCmd
}

func collectAvailableSubCmds(cmds []*ffcli.Command, args []string) error {
	for _, c := range cmds {
		subtitle := c.ShortHelp
		if f := c.FlagSet.Lookup(utils.AllNamespacesFlag); f != nil {
			subtitle = fmt.Sprintf("%s, opts [-%s: %s]", c.ShortHelp, f.Name, f.Usage)
		}

		if c.Name == versioncmd.CmdName {
			continue
		}
		if c.Name == portforwardcmd.CmdName {
			continue
		}

		rootConfig.Awf().Append(
			alfred.NewItem().
				SetTitle(c.Name).
				SetSubtitle(subtitle).
				SetAutocomplete(c.Name).
				SetVariable(utils.NextActionKey, utils.NextActionCmd).
				SetArg(c.Name),
		)
	}

	rootConfig.Awf().Filter(utils.GetQuery(args, 0)).Output()
	return nil
}
