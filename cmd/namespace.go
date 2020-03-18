package cmd

import (
	"fmt"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewNamespaceCmd create a new cmd for namespace resource
func NewNamespaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "namespace",
		Short: "list namespaces in current context",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			listNamespaces()
		},
		SilenceUsage: true,
	}

	return cmd
}

func listNamespaces() {
	namespaces, err := k.GetNamespaces()
	if err != nil {
		awf.Fatal("fatal error occurs", err.Error())
		return
	}
	for _, ns := range namespaces {
		awf.Append(&alfred.Item{
			Title:        fmt.Sprintf("current [%t] %s", ns.Current, ns.Name),
			Subtitle:     fmt.Sprintf("status [%s] age [%s]", ns.Status, ns.Age),
			Autocomplete: ns.Name,
			Arg:          ns.Name,
		})
	}

	awf.Output()
}
