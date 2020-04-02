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
		exitWith(err)
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
	rootCmd.AddCommand(NewCRDCmd())
	return rootCmd
}

// NewRootCmd create a new cmd for root
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Short: "available commands",
		Run: func(cmd *cobra.Command, args []string) {
			showAvailableSubCmds(cmd)
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

func showAvailableSubCmds(cmd *cobra.Command) {
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
		})
	}
	awf.Output()
}

func addAllNamespaceFlag(cmd *cobra.Command, all *bool) {
	cmd.PersistentFlags().BoolVarP(all, "all", "a", false, "in all namespaces")
}
