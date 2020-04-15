package cmd

import (
	"fmt"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewDeploymentCmd create a new cmd for deployment resource
func NewDeploymentCmd() *cobra.Command {
	var all bool
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "list deployments",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			listDeployments(all, getQuery(args, 0))
		},
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
	}
	addAllNamespaceFlag(cmd, &all)

	return cmd
}

func listDeployments(all bool, query string) {
	key := fmt.Sprintf("deploy-%t", all)
	if err := awf.Cache(key).MaxAge(cacheTime).LoadItems().Err(); err == nil {
		awf.Filter(query).Output()
		return
	}
	defer func() {
		awf.Cache(key).StoreItems().Workflow().Filter(query).Output()
	}()

	deps, err := k.GetDeployments(all)
	exitWith(err)
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
}
