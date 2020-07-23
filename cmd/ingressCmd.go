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
		title := getNamespacedResourceTitle(i)
		awf.Append(
			alfred.NewItem().
				SetTitle(title).
				SetSubtitle(
					fmt.Sprintf("host [%s] address [%s] ports [%s] ", i.Hosts, i.Address, i.Ports),
				).
				SetArg(i.Name).
				SetMod(alfred.ModCtrl,
					alfred.NewMod().
						SetSubtitle("copy ingress Address").
						SetArg(i.Address),
				),
		)
	}
	return
}
