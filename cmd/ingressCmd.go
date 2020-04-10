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
	key := fmt.Sprintf("ingress-%t", all)
	if err := awf.Cache(key).MaxAge(cacheTime).LoadItems().Err(); err == nil {
		awf.Filter(query).Output()
		return
	}
	defer func() {
		awf.Cache(key).StoreItems().Workflow().Filter(query).Output()
	}()

	ingresses, err := k.GetIngresses(all)
	exitWith(err)
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
}
