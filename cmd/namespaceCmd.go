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
		Short: "list namespaces in current context",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if ns == "" {
				listNamespaces()
				return nil
			}
			return setNamespace(ns)
		},
		SilenceUsage: true,
	}
	cmd.PersistentFlags().StringVarP(&ns, "name", "n", "", "namespace name")

	return cmd
}

func setNamespace(ns string) error {
	return k.SetNamespace(ns)
}

func listNamespaces() {
	namespaces, err := k.GetNamespaces()
	if err != nil {
		awf.Fatal("fatal error occurs", err.Error())
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
		})
	}

	awf.Output()
}
