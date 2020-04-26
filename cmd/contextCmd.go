package cmd

import (
	"fmt"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewContextCmd create a new cmd for context resource
func NewContextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "list contexts",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			use, del := getBoolFlag(cmd, useFalg), getBoolFlag(cmd, deleteFlag)
			if use {
				return shellOutputMiddleware(clearCacheMiddleware(useContext))(cmd, args)
			}
			if del {
				return shellOutputMiddleware(clearCacheMiddleware(deleteContext))(cmd, args)
			}
			return cacheOutputMiddleware(collectContexts)(cmd, args)
		},
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
	}
	addUseFlag(cmd)
	addDeleteFlag(cmd)

	return cmd
}

func useContext(cmd *cobra.Command, args []string) (err error) {
	context := getQuery(args, 0)
	if err = k.UseContext(context); err != nil {
		return
	}
	return
}

func deleteContext(cmd *cobra.Command, args []string) (err error) {
	context := getQuery(args, 0)
	if _, err = k.Execute(fmt.Sprintf("config delete-context %s", context)); err != nil {
		return
	}
	return
}

func collectContexts(cmd *cobra.Command, args []string) (err error) {
	contexts, err := k.GetContexts()
	if err != nil {
		return
	}

	for _, c := range contexts {
		title := c.Name
		if c.Current {
			title = fmt.Sprintf("[*] %s", c.Name)
		}

		// overwrite Arg for special case
		name := cmd.Name()
		deleteMod := getDeleteMod(name, c)
		deleteMod.Arg = fmt.Sprintf("%s %s --%s", name, c.Name, deleteFlag)
		awf.Append(&alfred.Item{
			Title: title,
			Arg:   c.Name,
			Mods: map[alfred.ModKey]alfred.Mod{
				alfred.ModCtrl:  getUseMod(name, c),
				alfred.ModShift: deleteMod,
			},
		})
	}
	return
}
