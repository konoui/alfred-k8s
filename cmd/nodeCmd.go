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
		Run: func(cmd *cobra.Command, args []string) {
			listNodes(getQuery(args, 0))
		},
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
	}
	return cmd
}

func listNodes(query string) {
	key := "node"
	if err := awf.Cache(key).MaxAge(cacheTime).LoadItems().Err(); err == nil {
		awf.Filter(query).Output()
		return
	}
	defer func() {
		awf.Cache(key).StoreItems().Workflow().Filter(query).Output()
	}()

	nodes, err := k.GetNodes()
	exitWith(err)
	for _, n := range nodes {
		awf.Append(&alfred.Item{
			Title:    n.Name,
			Subtitle: fmt.Sprintf("status [%s] version [%s]", n.Status, n.Version),
			Arg:      n.Name,
		})
	}
}
