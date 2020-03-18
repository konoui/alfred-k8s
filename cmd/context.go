package cmd

import (
	"fmt"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewContextCmd create a new cmd for context resource
func NewContextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "list context",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			listContexts()
		},
		SilenceUsage: true,
	}

	return cmd
}

func listContexts() {
	contexts, err := k.GetContexts()
	if err != nil {
		awf.Fatal("fatal error occurs", err.Error())
		return
	}
	for _, c := range contexts {
		awf.Append(&alfred.Item{
			Title:        fmt.Sprintf("current [%t] %s", c.Current, c.Name),
			Autocomplete: c.Name,
			Arg:          c.Name,
		})
	}

	awf.Output()
}
