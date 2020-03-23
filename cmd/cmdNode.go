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
			listNodes()
		},
		SilenceUsage: true,
	}
	return cmd
}

func listNodes() {
	nodes, err := k.GetNodes()
	if err != nil {
		awf.Fatal(fatalMessage, err.Error())
		return
	}
	for _, n := range nodes {

		awf.Append(&alfred.Item{
			Title:        n.Name,
			Subtitle:     fmt.Sprintf("status [%s] version [%s]", n.Status, n.Version),
			Autocomplete: n.Name,
			Arg:          n.Name,
		})
	}

	awf.Output()
}
