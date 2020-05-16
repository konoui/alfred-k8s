package cmd

import (
	"fmt"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewBaseCmd create a new cmd for resources not supported
func NewBaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "obj",
		Short: "list specific resources",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return outputMiddleware(collectBaseResources)(cmd, args)
		},
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
	}
	addAllNamespacesFlag(cmd)

	return cmd
}

func collectBaseResources(cmd *cobra.Command, args []string) (err error) {
	all := getBoolFlag(cmd, allNamespacesFlag)
	name := getQuery(args, 0)
	reses, err := k.GetBaseResources(name, all)
	if err != nil {
		return
	}

	for _, r := range reses {
		title := getNamespaceResourceTitle(r)
		awf.Append(&alfred.Item{
			Title:    title,
			Subtitle: fmt.Sprintf("age [%s]", r.Age),
			Arg:      r.Name,
		})
	}
	return
}
