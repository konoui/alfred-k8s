package cmd

import (
	"fmt"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewIngressCmd create a new cmd for ingress resource
func NewIngressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ingress",
		Short:   "list ingresses",
		Aliases: []string{"ing"},
		Args:    cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cacheOutputMiddleware(collectIngresses)(cmd, args)
		},
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
	}
	addAllNamespacesFlag(cmd)

	return cmd
}

func collectIngresses(cmd *cobra.Command, args []string) (err error) {
	all := getBoolFlag(cmd, allNamespacesFlag)
	ingresses, err := k.GetIngresses(all)
	if err != nil {
		return
	}
	for _, i := range ingresses {
		title := getNamespaceResourceTitle(i)
		awf.Append(&alfred.Item{
			Title:    title,
			Subtitle: fmt.Sprintf("host [%s] address [%s] ports [%s] ", i.Hosts, i.Address, i.Ports),
			Arg:      i.Name,
			Mods: map[alfred.ModKey]alfred.Mod{
				alfred.ModCtrl: {
					Subtitle: "copy ingress Address",
					Arg:      i.Address,
				},
			},
		})
	}
	return
}
