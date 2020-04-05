package cmd

import (
	"fmt"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewNamespaceCmd create a new cmd for namespace resource
func NewNamespaceCmd() *cobra.Command {
	var use bool
	cmd := &cobra.Command{
		Use:   "ns",
		Short: "list namespaces",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if use {
				useNamespace(getQuery(args, 0))
				return
			}
			listNamespaces(getQuery(args, 0))
		},
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
	}
	addUseFlag(cmd, &use)

	return cmd
}

func useNamespace(ns string) {
	if err := k.UseNamespace(ns); err != nil {
		fmt.Fprintf(errStream, "Failed due to %s\n", err)
		return
	}
	fmt.Fprintf(outStream, "Success!! switched %s namespace\n", ns)
}

func listNamespaces(query string) {
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
			Title:    title,
			Subtitle: fmt.Sprintf("status [%s] age [%s]", ns.Status, ns.Age),
			Arg:      ns.Name,
			Mods: map[alfred.ModKey]alfred.Mod{
				alfred.ModCtrl: alfred.Mod{
					Subtitle: "switch to the namespace",
					Arg:      fmt.Sprintf("ns %s --use", ns.Name),
					Variables: map[string]string{
						nextActionKey: nextActionShell,
					},
				},
			},
		})
	}

	awf.Filter(query).Output()
}
