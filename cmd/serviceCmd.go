package cmd

import (
	"fmt"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewServiceCmd create a new cmd for service resource
func NewServiceCmd() *cobra.Command {
	var all bool
	cmd := &cobra.Command{
		Use:   "svc",
		Short: "list services",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			listServices(all, getQuery(args, 0))
		},
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
	}
	addAllNamespaceFlag(cmd, &all)

	return cmd
}

func listServices(all bool, query string) {
	svcs, err := k.GetServices(all)
	if err != nil {
		awf.Fatal(fatalMessage, err.Error())
		return
	}

	for _, s := range svcs {
		title := getNamespaceResourceTitle(s)
		awf.Append(&alfred.Item{
			Title:    title,
			Subtitle: fmt.Sprintf("cluster-ip [%s] external-ip [%s] ports [%s]", s.ClusterIP, s.ExternalIP, s.Ports),
			Arg:      s.Name,
		})
	}

	awf.Filter(query).Output()
}
