package cmd

import (
	"fmt"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewPodCmd create a new cmd for pod resource
func NewPodCmd() *cobra.Command {
	var all, del bool
	cmd := &cobra.Command{
		Use:   "pod",
		Short: "list pods",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if del {
				deleteResource("pod", getQuery(args, 0), getQuery(args, 1))
				return
			}
			listPods(all, getQuery(args, 0))
		},
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
	}
	addDeleteFlag(cmd, &del)
	addAllNamespaceFlag(cmd, &all)

	return cmd
}

func listPods(all bool, query string) {
	key := fmt.Sprintf("pod-%t", all)
	if err := awf.Cache(key).MaxAge(cacheTime).LoadItems().Err(); err == nil {
		awf.Filter(query).Output()
		return
	}
	defer func() {
		awf.Cache(key).StoreItems().Workflow().Filter(query).Output()
	}()

	pods, err := k.GetPods(all)
	exitWith(err)
	for _, p := range pods {
		title := getNamespaceResourceTitle(p)
		awf.Append(&alfred.Item{
			Title:    title,
			Subtitle: fmt.Sprintf("ready [%s] status [%s] restarts [%s] ", p.Ready, p.Status, p.Restarts),
			Arg:      p.Name,
			Mods: map[alfred.ModKey]alfred.Mod{
				alfred.ModCtrl: {
					Subtitle: "rm the pod",
					Arg:      fmt.Sprintf("pod %s %s --delete", p.Name, p.Namespace),
					Variables: map[string]string{
						nextActionKey: nextActionShell,
					},
				},
				alfred.ModShift: getSternMod(p),
			},
		})
	}
}

func deleteResource(rs, name, ns string) {
	// TODO Clear cache
	arg := fmt.Sprintf("delete %s %s", rs, name)
	if ns != "" {
		arg = fmt.Sprintf("%s --namespace %s", arg, ns)
	}
	if _, err := k.Execute(arg); err != nil {
		fmt.Fprintf(outStream, "failed due to %s", err)
		return
	}
	fmt.Fprintf(outStream, "Success!! deleted the %s", rs)
}
