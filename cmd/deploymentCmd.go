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
	deps, err := k.GetDeployments(all)
	if err != nil {
		awf.Fatal(fatalMessage, err.Error())
		return
	}

	for _, d := range deps {
		title := d.Name
		if d.Namespace != "" {
			title = fmt.Sprintf("[%s] %s", d.Namespace, d.Name)
		}
		awf.Append(&alfred.Item{
			Title:    title,
			Subtitle: fmt.Sprintf("ready [%s] up-to-date [%s] available [%s]", d.Ready, d.UpToDate, d.Available),
			Arg:      d.Name,
		})
	}

	awf.Filter(query).Output()
}
