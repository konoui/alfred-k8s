package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/konoui/alfred-k8s/pkg/kubectl"
	"github.com/konoui/go-alfred"
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

// NewRootCmd create a new cmd for root
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "alfred-k8s <query>",
		Short: "operate k8s resources",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
		SilenceUsage: true,
	}

	return rootCmd
}

func run() {
	awf := alfred.NewWorkflow()
	// alfred script filter read from only stdout
	awf.SetStreams(outStream, outStream)
	awf.EmptyWarning("There are no resources", "No matching")
	pods, err := kubectl.GetAllPods()
	if err != nil {
		awf.Fatal("fatal error occurs", err.Error())
		return
	}

	for _, p := range pods {
		awf.Append(&alfred.Item{
			Title:        p.Name,
			Subtitle:     fmt.Sprintf("ready [%s] status [%s] restarts [%s] ", p.Ready, p.Status, p.Restarts),
			Autocomplete: p.Name,
			Arg:          p.Name,
		})
	}

	awf.Output()
}
