package cmd

import (
	"fmt"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewDeploymentCmd create a new cmd for deployment resource
func NewDeploymentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "list deployments",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cacheOutputMiddleware(collectDeployments)(cmd, args)
		},
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
	}
	addAllNamespacesFlag(cmd)

	return cmd
}

func collectDeployments(cmd *cobra.Command, args []string) (err error) {
	all := getBoolFlag(cmd, allNamespacesFlag)
	deps, err := k.GetDeployments(all)
	if err != nil {
		return
	}
	for _, d := range deps {
		title := getNamespaceResourceTitle(d)
		awf.Append(&alfred.Item{
			Title:    title,
			Subtitle: fmt.Sprintf("ready [%s] up-to-date [%s] available [%s]", d.Ready, d.UpToDate, d.Available),
			Arg:      d.Name,
			Mods: map[alfred.ModKey]alfred.Mod{
				alfred.ModShift: getSternMod(d),
			},
		})
	}
	return
}
