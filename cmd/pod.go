package cmd

import (
	"fmt"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewPodCmd create a new cmd for pod resource
func NewPodCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pod",
		Short: "list pods",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			listPods()
		},
		SilenceUsage: true,
	}

	return cmd
}

func listPods() {
	pods, err := k.GetAllPods()
	if err != nil {
		awf.Fatal("fatal error occurs", err.Error())
		return
	}
	for _, p := range pods {
		awf.Append(&alfred.Item{
			Title:        p.Name,
			Subtitle:     fmt.Sprintf("ready [%s] status [%s] restarts [%s] ", p.Ready, p.Status, p.Restarts),
			Autocomplete: p.Name,
			Arg:          p.Name,
		})
	}

	awf.Output()
}
