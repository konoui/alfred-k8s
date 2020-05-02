package cmd

import (
	"fmt"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewNamespaceCmd create a new cmd for namespace resource
func NewNamespaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ns",
		Short:   "list namespaces",
		Aliases: []string{"namespace"},
		Args:    cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			use := getBoolFlag(cmd, useFalg)
			if use {
				return shellOutputMiddleware(clearCacheMiddleware(useNamespace))(cmd, args)
			}
			return cacheOutputMiddleware(collectNamespaces)(cmd, args)
		},
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
	}
	addUseFlag(cmd)

	return cmd
}

func useNamespace(cmd *cobra.Command, args []string) (err error) {
	ns := getQuery(args, 0)
	if err = k.UseNamespace(ns); err != nil {
		return
	}
	return
}

func collectNamespaces(cmd *cobra.Command, args []string) (err error) {
	namespaces, err := k.GetNamespaces()
	if err != nil {
		return
	}
	for _, ns := range namespaces {
		title := ns.Name
		if ns.Current {
			title = fmt.Sprintf("[*] %s", ns.Name)
		}
		awf.Append(&alfred.Item{
			Title:    title,
			Subtitle: fmt.Sprintf("status [%s] age [%s]", ns.Status, ns.Age),
			Arg:      ns.Name,
			Mods: map[alfred.ModKey]alfred.Mod{
				alfred.ModCtrl: getUseMod("ns", ns),
			},
		})
	}
	return
}
