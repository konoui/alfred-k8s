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
	key := fmt.Sprintf("base-%t", all)
	if err := awf.Cache(key).MaxAge(cacheTime).LoadItems().Err(); err == nil {
		awf.Filter(query).Output()
		return
	}
	defer func() {
		awf.Cache(key).StoreItems().Workflow().Filter(query).Output()
	}()

	rs, err := k.GetBaseResources(name, all)
	exitWith(err)
	for _, r := range rs {
		title := getNamespaceResourceTitle(r)
		awf.Append(&alfred.Item{
			Title:    title,
			Subtitle: fmt.Sprintf("age [%s]", r.Age),
			Arg:      r.Name,
		})
	}
}
