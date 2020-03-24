package cmd

import (
	"fmt"

	"github.com/konoui/alfred-k8s/pkg/kubectl"
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
			listPods(all)
		},
		SilenceUsage: true,
	}
	cmd.PersistentFlags().BoolVarP(&all, "all", "a", false, "list pods in all namespaces")
	return cmd
}

func listPods(all bool) {
	var err error
	var pods []*kubectl.Pod
	if all {
		pods, err = k.GetAllPods()
	} else {
		pods, err = k.GetPods()
	}
	if err != nil {
		awf.Fatal(fatalMessage, err.Error())
		return
	}
	for _, p := range pods {
		title := p.Name
		if p.Namespace != "" {
			title = fmt.Sprintf("[%s] %s", p.Namespace, p.Name)
		}
		awf.Append(&alfred.Item{
			Title:    title,
			Subtitle: fmt.Sprintf("ready [%s] status [%s] restarts [%s] ", p.Ready, p.Status, p.Restarts),
			Arg:      p.Name,
		})
	}

	awf.Output()
}
