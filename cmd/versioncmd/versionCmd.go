package versioncmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/konoui/alfred-k8s/cmd/rootcmd"
	"github.com/peterbourgon/ff/v3/ffcli"
)

var (
	version  = "*"
	revision = "*"
)

const CmdName = "version"

// New create a new cmd for version
func New(rootConfig *rootcmd.Config) *ffcli.Command {
	fs := flag.NewFlagSet(CmdName, flag.ContinueOnError)

	cmd := &ffcli.Command{
		Name:      CmdName,
		ShortHelp: "print alfred-k8s version",
		FlagSet:   fs,
		Exec: func(ctx context.Context, args []string) error {
			fmt.Fprintf(rootConfig.Stdout(), "alfred-k8s %s (%s)\n", version, revision)
			return nil
		},
	}
	return cmd
}
