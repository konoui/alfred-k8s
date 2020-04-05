package cmd

import (
	"fmt"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewIngressCmd create a new cmd for ingress resource
func NewIngressCmd() *cobra.Command {
	var all bool
	cmd := &cobra.Command{
		Use:   "ingress",
		Short: "list ingresses",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			listIngresses(all, getQuery(args, 0))
		},
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
	}
	addAllNamespaceFlag(cmd, &all)

	return cmd
}

func listIngresses(all bool, query string) {
	ingresses, err := k.GetIngresses(all)
	if err != nil {
		awf.Fatal(fatalMessage, err.Error())
		return
	}

	for _, i := range ingresses {
		title := getNamespaceResourceTitle(i)
		awf.Append(&alfred.Item{
			Title:    title,
			Subtitle: fmt.Sprintf("host [%s] address [%s] ports [%s] ", i.Hosts, i.Address, i.Ports),
			Arg:      i.Name,
			Mods: map[alfred.ModKey]alfred.Mod{
				alfred.ModCtrl: alfred.Mod{
					Subtitle: "copy ingress Address",
					Arg:      i.Address,
				},
			},
		})
	}

	awf.Filter(query).Output()
}
