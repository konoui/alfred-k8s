package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/konoui/alfred-k8s/pkg/kubectl"
	"github.com/konoui/go-alfred"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// NewPortForwardCmd create a new cmd for port-forward
func NewPortForwardCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "port-forward <resource>/<resource-name>",
		Short: "list port-forwarded resources",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (non error) {
			use := getBoolFlag(cmd, useFalg)
			del := getBoolFlag(cmd, deleteFlag)
			if use {
				if err := background(cmd, args); err != nil {
					fmt.Fprintf(outStream, "Failed to execute port-forward due to %s", err)
				}
				return
			}

			if del {
				if err := terminateJob(cmd, args); err != nil {
					fmt.Fprintf(outStream, "Failed to terminate job due to %s. trying to scan jobs", err)
					_ = listPortResources()
				}
				fmt.Fprintln(outStream, "Success")
				return
			}

			_ = outputMiddleware(collectJobs)(cmd, args)
			return
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

func collectJobs(cmd *cobra.Command, args []string) error {
	jobs := listPortResources()
	for _, j := range jobs {
		awf.Append(&alfred.Item{
			Title:    fmt.Sprintf("[%s] %s/%s", j.Namespace, j.Resource, j.Name),
			Subtitle: fmt.Sprintf("ports %s", strings.Join(j.Ports, " ")),
			Mods: map[alfred.ModKey]alfred.Mod{
				alfred.ModAlt: {
					Subtitle: "terminate port-forward process",
					Arg:      fmt.Sprintf("%s %s/%s --%s %s --%s", cmd.Name(), j.Resource, j.Name, namespaceFlag, j.Namespace, deleteFlag),
					Variables: map[string]string{
						nextActionKey: nextActionShell,
					},
				},
			},
		})
	}
	return nil
}

// ListJobs return current jobs
func listPortResources() (jobs []*kubectl.PortResource) {
	dir := getDataDir()
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(errStream, "Invalid directory %s", dir)
		return nil
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if !strings.HasSuffix(f.Name(), ".pid") {
			continue
		}
		pidfile := filepath.Join(dir, f.Name())
		// try to read pidfile
		var j kubectl.PortResource
		err = readPidfile(pidfile, &j)
		if err != nil {
			fmt.Fprintln(errStream, err)
			continue
		}

		// valid process
		fmt.Fprintf(errStream, "Found the job %v\n", j)
		jobs = append(jobs, &j)
	}
	return jobs
}

// Note: terminateJob does not delete pid file.
// listJobs() takes responsibility for checking existing pid file, scaning process and deleting pidfile.
func terminateJob(cmd *cobra.Command, args []string) error {
	pidfile := getDataPath(cmd, args, ".pid")
	f, err := os.Open(pidfile)
	if err != nil {
		return err
	}
	var j kubectl.PortResource
	err = json.NewDecoder(f).Decode(&j)
	if err != nil {
		return err
	}

	err = j.Terminate()
	return errors.Wrapf(err, "pid %d", j.Pid)
}

func background(cmd *cobra.Command, args []string) error {
	res, name, ns := getResourceNameNamespace(cmd, args)
	ports := k.GetPorts(res, name, ns)
	if len(ports) == 0 {
		return fmt.Errorf("%s/%s has no ports", res, name)
	}

	// kill if the same command is executed
	pidfile := getDataPath(cmd, args, ".pid")
	prev := new(kubectl.PortResource)
	_ = readPidfile(pidfile, prev)
	_ = prev.Terminate()

	// create log file
	logfile := getDataPath(cmd, args, ".log")
	lf, err := os.Create(logfile)
	if err != nil {
		return errors.Wrapf(err, "failed to create log file %s", logfile)
	}
	defer lf.Close()
	// set logger
	jobStream = lf

	// do port forward
	job, errChan := k.PortForward(context.Background(), res, name, ns, ports)
	go func() {
		for line := range job.ReadLine() {
			fmt.Fprintln(jobStream, line)
		}
	}()

	f, err := os.Create(pidfile)
	if err != nil {
		return errors.Wrap(err, "failed to create pidfile")
	}
	defer f.Close()

	// write job information
	err = json.NewEncoder(f).Encode(job)
	if err != nil {
		return errors.Wrap(err, "failed to save data into pidfile")
	}

	// wait for port forward command
	return <-errChan
}

func readPidfile(pidfile string, j *kubectl.PortResource) error {
	v, err := ioutil.ReadFile(pidfile)
	logfile := strings.TrimRight(pidfile, "pid") + "log"
	if err != nil {
		_ = os.Remove(pidfile)
		_ = os.Remove(logfile)
		return errors.Wrapf(err, "failed to read %s", pidfile)
	}

	err = json.Unmarshal(v, j)
	if err != nil {
		_ = os.Remove(pidfile)
		_ = os.Remove(logfile)
		return errors.Wrapf(err, "failed to unmarshal %s", pidfile)
	}

	err = syscall.Kill(j.Pid, syscall.Signal(0))
	if err != nil {
		_ = os.Remove(pidfile)
		_ = os.Remove(logfile)
		return errors.Wrapf(err, "process (%d) does not exist", j.Pid)
	}

	return nil
}

func getDataPath(cmd *cobra.Command, args []string, extension string) string {
	res, name, ns := getResourceNameNamespace(cmd, args)
	file := cmd.Name() + "-" + res + "-" + ns + "-" + name + extension
	return filepath.Join(getDataDir(), file)
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
