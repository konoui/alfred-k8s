package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

// NewPortForwardCmd create a new cmd for port-forward
func NewPortForwardCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "port-forward <resource>/<resource-name>",
		Short: "list port-forwarded resources",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := awf.SetJobDir(getDataDir()); err != nil {
				return err
			}
			use := getBoolFlag(cmd, useFalg)
			del := getBoolFlag(cmd, deleteFlag)
			if use {
				return startPortForward(cmd, args)
			}
			if del {
				return stopPortForward(cmd, args)
			}
			return listJobs()
		},
		Hidden:             !experimental,
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
	}
	addUseFlag(cmd)
	addDeleteFlag(cmd)
	addNamespaceFlag(cmd)

	return cmd
}

func listJobs() error {
	jobs := awf.ListJobs()
	for _, job := range jobs {
		awf.Append(&alfred.Item{
			Title: job.Name(),
		})
	}
	awf.Output()
	return nil
}

func startPortForward(cmd *cobra.Command, args []string) error {
	awf.Append(&alfred.Item{
		Title: "Starting port forwarding",
	})
	awf.Job(getJobName(cmd, args)).Logging().
		StartWithExit(os.Args[0], os.Args[1:]...).
		Clear()

	res, name, ns := getResourceNameNamespace(cmd, args)
	ports := k.GetPorts(res, name, ns)
	if len(ports) == 0 {
		return fmt.Errorf("%s/%s has no ports", res, name)
	}

	kargs := append([]string{
		"port-forward",
		res + "/" + name,
		"--namespace",
		ns,
	}, ports...)
	resp, err := k.Execute(strings.Join(kargs, " "))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for l := range resp.Readline() {
		fmt.Println(l)
	}
	return nil
}

func stopPortForward(cmd *cobra.Command, args []string) error {
	return awf.Job(getJobName(cmd, args)).Terminate()
}

func getJobName(cmd *cobra.Command, args []string) string {
	res, name, ns := getResourceNameNamespace(cmd, args)
	return cmd.Name() + "-" + res + "-" + ns + "-" + name
}

func getDataDir() string {
	return "./data"
}

func getResourceNameNamespace(cmd *cobra.Command, args []string) (res, name, ns string) {
	query := getQuery(args, 0)
	tmp := strings.Split(query, "/")
	res = getQuery(tmp, 0)
	name = getQuery(tmp, 1)
	ns = getStringFlag(cmd, namespaceFlag)
	if ns == "" {
		var err error
		ns, err = k.GetCurrentNamespace()
		if err != nil {
			return res, name, "default"
		}
	}
	return
}
