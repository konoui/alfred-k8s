package cmd

import (
	"fmt"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewServiceCmd create a new cmd for service resource
func NewServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "svc",
		Short: "list services",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return outputMiddleware(collectServices)(cmd, args)
		},
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
	}
	addAllNamespacesFlag(cmd)

	return cmd
}

func collectServices(cmd *cobra.Command, args []string) (err error) {
	all := getBoolFlag(cmd, allNamespacesFlag)
	svcs, err := k.GetServices(all)
	if err != nil {
		return
	}
	for _, s := range svcs {
		title := getNamespaceResourceTitle(s)
		awf.Append(&alfred.Item{
			Title:    title,
			Subtitle: fmt.Sprintf("cluster-ip [%s] external-ip [%s] ports [%s]", s.ClusterIP, s.ExternalIP, s.Ports),
			Arg:      s.Name,
			Mods: map[alfred.ModKey]alfred.Mod{
				alfred.ModShift: getSternMod(s),
			},
		})
	}
	return
}
