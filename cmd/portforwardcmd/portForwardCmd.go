package portforwardcmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/konoui/alfred-k8s/cmd/rootcmd"
	"github.com/konoui/alfred-k8s/cmd/utils"
	"github.com/konoui/alfred-k8s/pkg/kubectl"
	"github.com/konoui/go-alfred"
	"github.com/peterbourgon/ff/v3/ffcli"
)

type Config struct {
	cmdName    string
	fs         *flag.FlagSet
	use        bool
	delete     bool
	resource   string
	namespace  string
	rootConfig *rootcmd.Config
}

const CmdName = "port-forward"

// New create a new cmd for port-forward
func New(rootConfig *rootcmd.Config) *ffcli.Command {
	fs := flag.NewFlagSet(CmdName, flag.ContinueOnError)
	cfg := &Config{
		rootConfig: rootConfig,
		fs:         fs,
	}
	cfg.registerFlags()

	cmd := &ffcli.Command{
		Name:       CmdName,
		ShortUsage: "port-forward <resource-name> --type <pod/svc/deploy>",
		ShortHelp:  "list port-forwarded resources",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			if err := cfg.rootConfig.Awf().SetJobDir(getDataDir()); err != nil {
				return err
			}
			if cfg.use {
				return cfg.startPortForward()
			}
			if cfg.delete {
				return cfg.stopPortForward()
			}
			return cfg.listJobs()
		},
	}

	return cmd
}

func (cfg *Config) registerFlags() {
	cfg.fs.BoolVar(&cfg.use, utils.UseFlag, false, "use it")
	cfg.fs.BoolVar(&cfg.delete, utils.DeleteFlag, false, "stop it")
	cfg.fs.StringVar(&cfg.resource, "type", "", "resource type e.g) svc, pod, deploy")
	cfg.fs.StringVar(&cfg.namespace, utils.NamespaceFlag, "", "resource namespace")
}

func (cfg *Config) listJobs() error {
	jobs := cfg.rootConfig.Awf().ListJobs()
	for _, job := range jobs {
		cfg.rootConfig.Awf().Append(
			alfred.NewItem().
				SetTitle(job.Name()),
		)
	}
	cfg.rootConfig.Awf().Output()
	return nil
}

func (cfg *Config) startPortForward() error {
	cfg.rootConfig.Awf().Append(&alfred.Item{
		Title: "Starting port forwarding",
	})
	cfg.rootConfig.Awf().Job(cfg.getJobName()).Logging().
		StartWithExit(os.Args[0], os.Args[1:]...).
		Clear()

	res, name, ns := cfg.resource, cfg.fs.Arg(0), cfg.namespace
	ports := cfg.rootConfig.Kubeclt().GetPorts(res, name, ns)
	if len(ports) == 0 {
		return fmt.Errorf("%s/%s has no ports", res, name)
	}

	kargs := append([]string{
		"port-forward",
		res + "/" + name,
		"--namespace",
		ns,
	}, ports...)
	resp, err := cfg.rootConfig.Kubeclt().Execute(strings.Join(kargs, " "))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for l := range resp.Readline() {
		fmt.Println(l)
	}

	return nil
}

func (cfg *Config) stopPortForward() error {
	return cfg.rootConfig.Awf().Job(cfg.getJobName()).Terminate()
}

func (cfg *Config) getJobName() string {
	res, name, ns := cfg.resource, cfg.fs.Arg(0), cfg.namespace
	return cfg.cmdName + "-" + res + "-" + name + "-" + ns
}

func getDataDir() string {
	return "./data"
}

func GetCopyPortForwardMod(k *kubectl.Kubectl, res string, i interface{}) *alfred.Mod {
	name, ns := kubectl.GetNameNamespace(i)
	if ns == "" {
		var err error
		ns, err = k.GetCurrentNamespace()
		if err != nil {
			ns = "default"
		}
	}
	ports := k.GetPorts(res, name, ns)
	if len(ports) == 0 {
		return alfred.NewMod().
			SetSubtitle("the resource has no ports")
	}

	arg := fmt.Sprintf("kubectl port-forward %s/%s %s", res, name, strings.Join(ports, " "))
	if ns != "" {
		arg = fmt.Sprintf("%s --namespace %s", arg, ns)
	}

	return alfred.NewMod().
		SetSubtitle("copy " + arg).
		SetArg(arg)
}

func GetExecPortForwardMod(k *kubectl.Kubectl, res string, i interface{}) *alfred.Mod {
	name, ns := kubectl.GetNameNamespace(i)
	arg := fmt.Sprintf("--type %s --%s %s", res, utils.UseFlag, name)
	if ns != "" {
		arg = fmt.Sprintf("--%s %s %s", utils.NamespaceFlag, ns, arg)
	}
	cmd := CmdName + " " + arg

	return alfred.NewMod().
		SetSubtitle("port-forward in background").
		SetArg(cmd).
		SetVariable(utils.NextActionKey, utils.NextActionJob)
}
