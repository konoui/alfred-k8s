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
			if len(args) == 0 {
				listCustomResources(getQuery(args, 0))
				return
			}
			listSpecificResources(args[0], all, getQuery(args, 1))
		},
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
	}
	addAllNamespaceFlag(cmd, &all)

	return cmd
}

func listCustomResources(query string) {
	crds, err := k.GetCRDs()
	if err != nil {
		awf.Fatal(fatalMessage, err.Error())
		return
	}
	for _, c := range crds {
		awf.Append(&alfred.Item{
			Title:    c.Name,
			Subtitle: fmt.Sprintf("created-at [%s] ", c.CreatedAT),
			Arg:      c.Name,
		})
	}

	awf.Filter(query).Output()
}

func listSpecificResources(name string, all bool, query string) {
	rs, err := k.GetSpecificResources(name, all)
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
