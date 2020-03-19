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
		Short: "list context",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if context == "" {
				listContexts()
				return nil
			}
			return setContext(context)
		},
		SilenceUsage: true,
	}

	cmd.PersistentFlags().StringVarP(&context, "name", "n", "", "context name")
	return cmd
}

func setContext(context string) error {
	return k.SetContext(context)

}

func listContexts() {
	contexts, err := k.GetContexts()
	if err != nil {
		awf.Fatal("fatal error occurs", err.Error())
		return
	}
	for _, c := range contexts {
		title := c.Name
		if c.Current {
			title = fmt.Sprintf("[*] %s", c.Name)
		}
		awf.Append(&alfred.Item{
			Title:        title,
			Autocomplete: c.Name,
			Arg:          c.Name,
		})
	}

	awf.Output()
}
