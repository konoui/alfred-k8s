package cmd

import (
	"fmt"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewContextCmd create a new cmd for context resource
func NewContextCmd() *cobra.Command {
	var use bool
	cmd := &cobra.Command{
		Use:   "context",
		Short: "list contexts",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if use {
				useContext(getQuery(args, 0))
				return
			}
			listContexts(getQuery(args, 0))
		},
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
	}
	addUseFlag(cmd, &use)

	return cmd
}

func useContext(context string) {
	if err := k.UseContext(context); err != nil {
		fmt.Fprintf(errStream, "Failed due to %s\n", err)
		return
	}
	fmt.Fprintln(outStream, "Success!! switched context")
}

func listContexts(query string) {
	contexts, err := k.GetContexts()
	if err != nil {
		awf.Fatal(fatalMessage, err.Error())
		return
	}
	for _, c := range contexts {
		title := c.Name
		if c.Current {
			title = fmt.Sprintf("[*] %s", c.Name)
		}

		awf.Append(&alfred.Item{
			Title: title,
			Arg:   c.Name,
			Mods: map[alfred.ModKey]alfred.Mod{
				alfred.ModCtrl: alfred.Mod{
					Subtitle: "switch to the context",
					Arg:      fmt.Sprintf("context %s --use", c.Name),
					Variables: map[string]string{
						nextActionKey: nextActionShell,
					},
				},
			},
		})
	}

	awf.Filter(query).Output()
}
