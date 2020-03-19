package cmd

import (
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	outStream io.Writer = os.Stdout
	errStream io.Writer = os.Stderr
)

// Execute root cmd
func Execute(rootCmd *cobra.Command) {
	// Note: result of RunE redirects etderr
	rootCmd.SetOutput(errStream)
	if err := rootCmd.Execute(); err != nil {
		log.Printf("command execution failed: %+v", err)
		os.Exit(1)
	}
}

// NewCmd create sub commands
func NewCmd() *cobra.Command {
	rootCmd := NewRootCmd()
	rootCmd.AddCommand(NewPodCmd())
	rootCmd.AddCommand(NewContextCmd())
	rootCmd.AddCommand(NewNamespaceCmd())
	rootCmd.AddCommand(NewDeploymentCmd())
	rootCmd.AddCommand(NewServiceCmd())
	return rootCmd
}

// NewRootCmd create a new cmd for root
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{}
	return rootCmd
}
