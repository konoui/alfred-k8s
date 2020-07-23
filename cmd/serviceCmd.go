package cmd

import (
	"fmt"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewServiceCmd create a new cmd for service resource
func NewServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "svc",
		Short:   "list services",
		Aliases: []string{"service"},
		Args:    cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cacheOutputMiddleware(collectServices)(cmd, args)
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
		title := getNamespacedResourceTitle(s)
		awf.Append(
			alfred.NewItem().
				SetTitle(title).
				SetSubtitle(
					fmt.Sprintf("cluster-ip [%s] external-ip [%s] ports [%s]", s.ClusterIP, s.ExternalIP, s.Ports),
				).
				SetArg(s.Name).
				SetMod(alfred.ModShift, getSternMod(s)).
				SetMod(alfred.ModAlt, getPortForwardMod(cmd.Name(), s)),
		)
	}
	return
}
