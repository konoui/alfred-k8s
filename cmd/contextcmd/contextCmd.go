package contextcmd

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
	use        bool
	delete     bool
	fs         *flag.FlagSet
	rootConfig *rootcmd.Config
}

const CmdName = "context"

// New create a new cmd for context resource
func New(rootConfig *rootcmd.Config) *ffcli.Command {
	fs := flag.NewFlagSet(CmdName, flag.ContinueOnError)
	cfg := Config{
		rootConfig: rootConfig,
		fs:         fs,
	}
	cfg.registerFlags()

	cmd := &ffcli.Command{
		Name:      CmdName,
		ShortHelp: "list contexts",
		FlagSet:   fs,
		Exec: func(ctx context.Context, args []string) error {
			if cfg.delete {
				return cfg.deleteContext()
			}
			if cfg.use {
				return cfg.useContext()
			}
			return cfg.collectContexts()
		},
	}

	return cmd
}

func (cfg *Config) registerFlags() {
	cfg.fs.BoolVar(&cfg.use, utils.UseFlag, false, "use it")
	cfg.fs.BoolVar(&cfg.delete, utils.DeleteFlag, false, "delete it")
}

func (cfg *Config) useContext() (err error) {
	contextName := cfg.fs.Arg(0)
	if err = cfg.rootConfig.Kubeclt().UseContext(contextName); err != nil {
		fmt.Fprintf(cfg.rootConfig.Stdout(), "Failed due to %s\n", err)
		return nil
	}
	fmt.Fprintf(cfg.rootConfig.Stdout(), "Success!!\n")
	return
}

func (cfg *Config) deleteContext() (err error) {
	contextName := cfg.fs.Arg(0)
	if _, err = cfg.rootConfig.Kubeclt().Execute(fmt.Sprintf("config delete-context %s", contextName)); err != nil {
		fmt.Fprintf(cfg.rootConfig.Stdout(), "Failed due to %s\n", err)
		return nil
	}
	fmt.Fprintf(cfg.rootConfig.Stdout(), "Success!!\n")
	return
}

func (cfg *Config) collectContexts() (err error) {
	contexts, err := cfg.rootConfig.Kubeclt().GetContexts()
	if err != nil {
		return
	}

	for _, c := range contexts {
		title := c.Name
		if c.Current {
			title = fmt.Sprintf("[*] %s", c.Name)
		}

		// overwrite Arg for special case as context is non namespaced resource but `c` has namespace field.
		deleteMod := utils.GetDeleteMod(CmdName, c)
		deleteMod.Arg = fmt.Sprintf("%s --%s %s", CmdName, utils.DeleteFlag, c.Name)
		useMod := utils.GetUseMod(CmdName, c)
		useMod.Arg = fmt.Sprintf("%s --%s %s", CmdName, utils.UseFlag, c.Name)
		cfg.rootConfig.Awf().Append(
			alfred.NewItem().
				SetTitle(title).
				SetArg(c.Name).
				SetMod(alfred.ModCtrl, useMod).
				SetMod(alfred.ModShift, deleteMod),
		)
	}

	cfg.rootConfig.Awf().Filter(cfg.fs.Arg(0)).Output()
	return
}
