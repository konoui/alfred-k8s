package cmd

import (
	"fmt"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewPodCmd create a new cmd for pod resource
func NewPodCmd() *cobra.Command {
	var all bool
	cmd := &cobra.Command{
		Use:   "pod",
		Short: "list pods",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			listPods(all, getQuery(args, 0))
		},
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
	}
	addAllNamespaceFlag(cmd, &all)

	return cmd
}

func listPods(all bool, query string) {
	pods, err := k.GetPods(all)
	if err != nil {
		awf.Fatal(fatalMessage, err.Error())
		return
	}

	for _, p := range pods {
		title := getNamespaceResourceTitle(p)
		awf.Append(&alfred.Item{
			Title:    title,
			Subtitle: fmt.Sprintf("ready [%s] status [%s] restarts [%s] ", p.Ready, p.Status, p.Restarts),
			Arg:      p.Name,
		})
	}

	awf.Filter(query).Output()
}
