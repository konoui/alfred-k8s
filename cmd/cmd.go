package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

var (
	outStream io.Writer = os.Stdout
	errStream io.Writer = os.Stderr
)

// Execute root cmd
func Execute(rootCmd *cobra.Command) {
	rootCmd.SetErr(errStream)
	rootCmd.SetOut(outStream)
	if err := rootCmd.Execute(); err != nil {
		listAvailableSubCmds(rootCmd, getQuery(os.Args, 1))
	}
}

// NewDefaultCmd create sub commands
func NewDefaultCmd() *cobra.Command {
	rootCmd := NewRootCmd()
	rootCmd.AddCommand(NewPodCmd())
	rootCmd.AddCommand(NewContextCmd())
	rootCmd.AddCommand(NewNamespaceCmd())
	rootCmd.AddCommand(NewDeploymentCmd())
	rootCmd.AddCommand(NewServiceCmd())
	rootCmd.AddCommand(NewNodeCmd())
	rootCmd.AddCommand(NewIngressCmd())
	rootCmd.AddCommand(NewBaseCmd())
	rootCmd.AddCommand(NewCRDCmd())
	return rootCmd
}

// NewRootCmd create a new cmd for root
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Short: "list available commands",
		Run: func(cmd *cobra.Command, args []string) {
			listAvailableSubCmds(cmd, getQuery(args, 0))
		},
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
	}
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
	return rootCmd
}

func listAvailableSubCmds(cmd *cobra.Command, query string) {
	defer func() {
		awf.Filter(query).Output()
	}()

	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() {
			continue
		}

		subtitle := c.Short
		if f := c.Flag("all"); f != nil {
			subtitle = fmt.Sprintf("%s, opts [-%s: %s]", c.Short, f.Shorthand, f.Usage)
		}

		awf.Append(&alfred.Item{
			Title:        c.Name(),
			Subtitle:     subtitle,
			Autocomplete: c.Name(),
			Variables: map[string]string{
				nextActionKey: nextActionCmd,
			},
			Arg: c.Name(),
		})
	}
}
