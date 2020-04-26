package cmd

import (
	"fmt"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewNodeCmd create a new cmd for node resource
func NewNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node",
		Short: "list nodes",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cacheOutputMiddleware(collectNodes)(cmd, args)
		},
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
	}
	return cmd
}

func collectNodes(cmd *cobra.Command, args []string) (err error) {
	nodes, err := k.GetNodes()
	if err != nil {
		return
	}
	for _, n := range nodes {
		awf.Append(&alfred.Item{
			Title:    n.Name,
			Subtitle: fmt.Sprintf("status [%s] version [%s]", n.Status, n.Version),
			Arg:      n.Name,
		})
	}
	return
}
