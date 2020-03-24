package cmd

import (
	"fmt"

	"github.com/konoui/alfred-k8s/pkg/kubectl"
	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewServiceCmd create a new cmd for service resource
func NewServiceCmd() *cobra.Command {
	var all bool
	cmd := &cobra.Command{
		Use:   "svc",
		Short: "list services",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			listServices(all)
		},
		SilenceUsage: true,
	}
	cmd.PersistentFlags().BoolVarP(&all, "all", "a", false, "list services in all namespaces")
	return cmd
}

func listServices(all bool) {
	var err error
	var svcs []*kubectl.Service
	if all {
		svcs, err = k.GetAllServices()
	} else {
		svcs, err = k.GetServices()
	}
	if err != nil {
		awf.Fatal(fatalMessage, err.Error())
		return
	}
	for _, s := range svcs {
		title := fmt.Sprintf("type [%s]  %s", s.Type, s.Name)
		if s.Namespace != "" {
			title = fmt.Sprintf("[%s] %s", s.Namespace, s.Name)
		}
		awf.Append(&alfred.Item{
			Title:    title,
			Subtitle: fmt.Sprintf("cluster-ip [%s] external-ip [%s] ports [%s]", s.ClusterIP, s.ExternalIP, s.Ports),
			Arg:      s.Name,
		})
	}

	awf.Output()
}
