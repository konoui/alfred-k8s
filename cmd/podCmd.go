package cmd

import (
	"fmt"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewPodCmd create a new cmd for pod resource
func NewPodCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pod",
		Short:   "list pods",
		Aliases: []string{"po"},
		Args:    cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			del := getBoolFlag(cmd, deleteFlag)
			if del {
				return shellOutputMiddleware(clearCacheMiddleware(deleteResource))(cmd, args)
			}
			return cacheOutputMiddleware(collectPods)(cmd, args)
		},
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
	}
	addDeleteFlag(cmd)
	addNamespaceFlag(cmd)
	addAllNamespacesFlag(cmd)

	return cmd
}

func collectPods(cmd *cobra.Command, args []string) error {
	all := getBoolFlag(cmd, allNamespacesFlag)
	pods, err := k.GetPods(all)
	if err != nil {
		return err
	}
	for _, p := range pods {
		title := getNamespaceResourceTitle(p)
		awf.Append(&alfred.Item{
			Title:    title,
			Subtitle: fmt.Sprintf("ready [%s] status [%s] restarts [%s] ", p.Ready, p.Status, p.Restarts),
			Arg:      p.Name,
			Mods: map[alfred.ModKey]alfred.Mod{
				alfred.ModCtrl:  getDeleteMod(cmd.Name(), p),
				alfred.ModShift: getSternMod(p),
				alfred.ModAlt:   getPortForwardMod(cmd.Name(), p),
			},
		})
	}
	return nil
}

func deleteResource(cmd *cobra.Command, args []string) (err error) {
	// resource name must be same as cobra.Command Use
	rs := cmd.Name()
	name := getQuery(args, 0)
	ns := getStringFlag(cmd, namespaceFlag)

	arg := fmt.Sprintf("delete %s %s", rs, name)
	if ns != "" {
		arg = fmt.Sprintf("%s --namespace %s", arg, ns)
	}
	if _, err := k.Execute(arg); err != nil {
		return err
	}
	return nil
}
