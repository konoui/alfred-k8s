package cmd

import (
	"fmt"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewNamespaceCmd create a new cmd for namespace resource
func NewNamespaceCmd() *cobra.Command {
	var ns string
	cmd := &cobra.Command{
		Use:   "ns",
		Short: "list namespaces",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if ns == "" {
				listNamespaces()
				return
			}
			useNamespace(ns)
		},
		SilenceUsage: true,
	}
	cmd.PersistentFlags().StringVarP(&ns, "name", "n", "", "namespace name to switch")

	return cmd
}

func useNamespace(ns string) {
	if err := k.UseNamespace(ns); err != nil {
		fmt.Fprintf(errStream, "Failed due to %s\n", err)
		return
	}
	fmt.Fprintf(outStream, "Success!! switched %s namespace\n", ns)
}

func listNamespaces() {
	namespaces, err := k.GetNamespaces()
	if err != nil {
		awf.Fatal(fatalMessage, err.Error())
		return
	}
	for _, ns := range namespaces {
		title := ns.Name
		if ns.Current {
			title = fmt.Sprintf("[*] %s", ns.Name)
		}
		awf.Append(&alfred.Item{
			Title:        title,
			Subtitle:     fmt.Sprintf("status [%s] age [%s]", ns.Status, ns.Age),
			Autocomplete: ns.Name,
			Arg:          ns.Name,
			Mods: map[alfred.ModKey]alfred.Mod{
				alfred.ModCtrl: alfred.Mod{
					Subtitle: "switch to specific namespace",
					Arg:      fmt.Sprintf("ns --name %s", ns.Name),
					Variables: map[string]string{
						nextActionKey: nextActionCmd,
					},
				},
			},
		})
	}

	awf.Output()
}
