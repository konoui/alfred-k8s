package cmd

import (
	"fmt"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewContextCmd create a new cmd for context resource
func NewContextCmd() *cobra.Command {
	var context string
	cmd := &cobra.Command{
		Use:   "context",
		Short: "list contexts",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if context == "" {
				listContexts()
				return
			}
			useContext(context)
		},
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
	}

	cmd.PersistentFlags().StringVarP(&context, "use", "u", "", "context name to switch")
	return cmd
}

func useContext(context string) {
	if err := k.UseContext(context); err != nil {
		fmt.Fprintf(errStream, "Failed due to %s\n", err)
		return
	}
	fmt.Fprintln(outStream, "Success!! switched context")
}

func listContexts() {
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
					Subtitle: "switch to specific context",
					Arg:      fmt.Sprintf("context --use %s", c.Name),
					Variables: map[string]string{
						nextActionKey: nextActionCmd,
					},
				},
			},
		})
	}

	awf.Output()
}
