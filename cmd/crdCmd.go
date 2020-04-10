package cmd

import (
	"fmt"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewCRDCmd create a new cmd for crd resource
func NewCRDCmd() *cobra.Command {
	var all bool
	cmd := &cobra.Command{
		Use:   "crd",
		Short: "list crds",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			listCustomResources(getQuery(args, 0))
		},
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
	}
	addAllNamespaceFlag(cmd, &all)

	return cmd
}

func listCustomResources(query string) {
	key := "crd"
	if err := awf.Cache(key).MaxAge(cacheTime).LoadItems().Err(); err == nil {
		awf.Filter(query).Output()
		return
	}
	defer func() {
		awf.Cache(key).StoreItems().Workflow().Filter(query).Output()
	}()

	crds, err := k.GetCRDs()
	exitWith(err)
	for _, c := range crds {
		awf.Append(&alfred.Item{
			Title:    c.Name,
			Subtitle: fmt.Sprintf("created-at [%s] ", c.CreatedAT),
			Arg:      c.Name,
		})
	}
}
