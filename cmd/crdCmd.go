package cmd

import (
	"fmt"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewCRDCmd create a new cmd for crd resource
func NewCRDCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "crd",
		Short: "list crds",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cacheOutputMiddleware(collectCRDs)(cmd, args)
		},
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
	}

	return cmd
}

func collectCRDs(cmd *cobra.Command, args []string) (err error) {
	crds, err := k.GetCRDs()
	if err != nil {
		return
	}
	for _, c := range crds {
		awf.Append(&alfred.Item{
			Title:    c.Name,
			Subtitle: fmt.Sprintf("created-at [%s] ", c.CreatedAT),
			Arg:      c.Name,
		})
	}
	return
}
