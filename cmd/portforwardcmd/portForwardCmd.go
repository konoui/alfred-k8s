package portforwardcmd

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/konoui/alfred-k8s/cmd/rootcmd"
	"github.com/konoui/alfred-k8s/cmd/utils"
	"github.com/konoui/alfred-k8s/pkg/kubectl"
	"github.com/konoui/go-alfred"
	"github.com/peterbourgon/ff/v3/ffcli"
)

type Config struct {
	fs           *flag.FlagSet
	use          bool
	delete       bool
	resourceType string
	name         string
	namespace    string
	rootConfig   *rootcmd.Config
}

const CmdName = "port-forward"
const typeFlag = "type"

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
			// TODO
			// Note set resource name after flag.Parse
			cfg.name = cfg.fs.Arg(0)
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
	cfg.fs.StringVar(&cfg.resourceType, typeFlag, "", "resource type e.g) svc, pod, deploy")
	cfg.fs.StringVar(&cfg.namespace, utils.NamespaceFlag, "", "resource namespace")
}

func (cfg *Config) listJobs() error {
	jobs := cfg.rootConfig.Awf().ListJobs()
	for _, job := range jobs {
		resType, name, ns := cfg.getValuesFromJobName(job.Name())
		cfg.rootConfig.Awf().Append(
			alfred.NewItem().
				SetTitle(job.Name()).
				SetMod(
					alfred.ModCtrl,
					alfred.NewMod().
						SetSubtitle("stop port forward").
						SetArg(
							fmt.Sprintf("%s --%s %s --%s %s --%s %s", CmdName, typeFlag, resType, utils.NamespaceFlag, ns, utils.DeleteFlag, name),
						).
						SetVariable(utils.NextActionKey, utils.NextActionShell),
				),
		)
	}
	cfg.rootConfig.Awf().Output()
	return nil
}

func (cfg *Config) startPortForward() error {
	cfg.rootConfig.Awf().Append(
		alfred.NewItem().
			SetTitle("Starting port forward"),
	)

	resType, name, ns := cfg.resourceType, cfg.name, cfg.namespace
	ports := cfg.rootConfig.Kubeclt().GetPorts(resType, name, ns)
	if len(ports) == 0 {
		return fmt.Errorf("%s/%s has no ports", resType, name)
	}

	kargs := append([]string{
		"port-forward",
		resType + "/" + name,
		"--namespace",
		ns,
	}, ports...)

	cmds, _ := cfg.rootConfig.Kubeclt().GetKubectlCommandEnv(kargs)
	cfg.rootConfig.Awf().Job(cfg.getJobName()).Logging().
		StartWithExit(cmds[0], cmds[1:]...).
		Clear()

	return nil
}

func (cfg *Config) stopPortForward() error {
	return cfg.rootConfig.Awf().Job(cfg.getJobName()).Terminate()
}

func (cfg *Config) getJobName() string {
	res, name, ns := cfg.resourceType, cfg.name, cfg.namespace
	return CmdName + "_" + res + "_" + name + "_" + ns
}

func (cfg *Config) getValuesFromJobName(jobName string) (res, name, namespace string) {
	values := strings.SplitN(jobName, "_", 4)
	if len(values) != 4 {
		return
	}
	res = values[1]
	name = values[2]
	namespace = values[3]
	cfg.resourceType = res
	cfg.name = name
	cfg.namespace = namespace
	return
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
