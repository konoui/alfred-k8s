package podcmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/konoui/alfred-k8s/cmd/rootcmd"
	"github.com/konoui/alfred-k8s/cmd/utils"
	"github.com/konoui/go-alfred"
	"github.com/peterbourgon/ff/v3/ffcli"
)

type Config struct {
	all        bool
	delete     bool
	namespace  string
	fs         *flag.FlagSet
	rootConfig *rootcmd.Config
}

const CmdName = "pod"

// New create a new cmd for pod resource
func New(rootConfig *rootcmd.Config) *ffcli.Command {
	fs := flag.NewFlagSet(CmdName, flag.ContinueOnError)
	cfg := &Config{
		rootConfig: rootConfig,
		fs:         fs,
	}
	cfg.registerFlags()

	cmd := &ffcli.Command{
		Name:      CmdName,
		ShortHelp: "list pods",
		FlagSet:   fs,
		Exec: func(ctx context.Context, args []string) error {
			if cfg.delete {
				return cfg.deleteResource()
			}
			return cfg.collectPods()
		},
	}

	return cmd
}

func (cfg *Config) registerFlags() {
	cfg.fs.BoolVar(&cfg.delete, utils.DeleteFlag, false, "delete it")
	cfg.fs.BoolVar(&cfg.all, utils.AllNamespacesFlag, false, "in all namespaces")

}

func (cfg *Config) collectPods() error {
	pods, err := cfg.rootConfig.Kubeclt().GetPods(cfg.all)
	if err != nil {
		return err
	}
	for _, p := range pods {
		title := utils.GetNamespacedResourceTitle(p)
		cfg.rootConfig.Awf().Append(
			alfred.NewItem().
				SetTitle(title).
				SetSubtitle(
					fmt.Sprintf("ready [%s] status [%s] restarts [%s] ", p.Ready, p.Status, p.Restarts),
				).
				SetArg(p.Name).
				SetMod(alfred.ModCtrl, utils.GetDeleteMod(CmdName, p)).
				SetMod(alfred.ModShift, utils.GetSternMod(p)).
				SetMod(alfred.ModAlt, utils.GetPortForwardMod(cfg.rootConfig.Kubeclt(), CmdName, p)),
		)
	}

	cfg.rootConfig.Awf().Filter(cfg.fs.Arg(0)).Output()
	return nil
}

func (cfg *Config) deleteResource() error {
	pod := cfg.fs.Arg(0)
	arg := fmt.Sprintf("delete %s %s", CmdName, pod)
	if cfg.namespace != "" {
		arg = fmt.Sprintf("%s --namespace %s", arg, cfg.namespace)
	}
	if _, err := cfg.rootConfig.Kubeclt().Execute(arg); err != nil {
		fmt.Fprintf(cfg.rootConfig.Stdout(), "Failed due to %s\n", err)
		return nil
	}

	fmt.Fprintf(cfg.rootConfig.Stdout(), "Success!!\n")
	return nil
}
