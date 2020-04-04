package cmd

import (
	"fmt"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewBaseCmd create a new cmd for resources not supported
func NewBaseCmd() *cobra.Command {
	var all bool
	cmd := &cobra.Command{
		Use:   "obj",
		Short: "list specific resources",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			listBaseResources(getQuery(args, 0), all, getQuery(args, 1))
		},
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
	}
	addAllNamespaceFlag(cmd, &all)

	return cmd
}

func listBaseResources(name string, all bool, query string) {
	rs, err := k.GetBaseResources(name, all)
	if err != nil {
		awf.Fatal(fatalMessage, err.Error())
		return
	}

	for _, r := range rs {
		title := r.Name
		if r.Namespace != "" {
			title = fmt.Sprintf("[%s] %s", r.Namespace, r.Name)
		}
		awf.Append(&alfred.Item{
			Title:    title,
			Subtitle: fmt.Sprintf("age [%s]", r.Age),
			Arg:      r.Name,
		})
	}

	awf.Filter(query).Output()
}
