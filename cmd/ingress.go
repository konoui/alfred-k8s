package cmd

import (
	"fmt"

	"github.com/konoui/alfred-k8s/pkg/kubectl"
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
			listIngresses(all)
		},
		SilenceUsage: true,
	}
	cmd.PersistentFlags().BoolVarP(&all, "all", "a", false, "list ingresses in all namespaces")
	return cmd
}

func listIngresses(all bool) {
	var err error
	var ingresses []*kubectl.Ingress
	if all {
		ingresses, err = k.GetAllIngresses()
	} else {
		ingresses, err = k.GetIngresses()
	}
	if err != nil {
		awf.Fatal(fatalMessage, err.Error())
		return
	}
	for _, i := range ingresses {
		title := i.Name
		if i.Namespace != "" {
			title = fmt.Sprintf("[%s] %s", i.Namespace, i.Name)
		}
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

	awf.Output()
}
