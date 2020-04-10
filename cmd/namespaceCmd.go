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
	key := "ns"
	if err := awf.Cache(key).MaxAge(cacheTime).LoadItems().Err(); err == nil {
		awf.Filter(query).Output()
		return
	}
	defer func() {
		awf.Cache(key).StoreItems().Workflow().Filter(query).Output()
	}()

	namespaces, err := k.GetNamespaces()
	exitWith(err)
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
				alfred.ModCtrl: {
					Subtitle: "switch to the namespace",
					Arg:      fmt.Sprintf("ns %s --use", ns.Name),
					Variables: map[string]string{
						nextActionKey: nextActionShell,
					},
				},
			},
		})
	}
}
