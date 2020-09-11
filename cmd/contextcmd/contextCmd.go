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
	cfg := &Config{
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
				return cfg.rootConfig.DeleteOutput(
					cfg,
					cfg.GetQuery(),
				)
			}
			if cfg.use {
				return cfg.rootConfig.UseOutput(
					cfg,
					cfg.GetQuery(),
				)
			}

			return cfg.rootConfig.CollectOutput(
				cfg,
				cfg.GetQuery(),
				utils.GetCacheKey(CmdName, false),
			)
		},
	}

	return cmd
}

func (cfg *Config) registerFlags() {
	cfg.fs.BoolVar(&cfg.use, utils.UseFlag, false, "use it")
	cfg.fs.BoolVar(&cfg.delete, utils.DeleteFlag, false, "delete it")
}

func (cfg *Config) Use() (err error) {
	contextName := cfg.fs.Arg(0)
	_ = cfg.rootConfig.Kubeclt().UseContext(contextName)
	return nil
}

func (cfg *Config) Delete() (err error) {
	contextName := cfg.fs.Arg(0)
	_, _ = cfg.rootConfig.Kubeclt().Execute(fmt.Sprintf("config delete-context %s", contextName))
	return nil
}

func (cfg *Config) Collect() error {
	contexts, err := cfg.rootConfig.Kubeclt().GetContexts()
	if err != nil {
		return err
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
		copyMod := utils.GetCopyMod("", c.Name)

		modMap := map[rootcmd.KeyMapKey]*alfred.Mod{
			rootcmd.CopyResourceKey:   copyMod,
			rootcmd.UseResourceKey:    useMod,
			rootcmd.DeleteResourceKey: deleteMod,
		}
		enterMod, mods := rootcmd.MakeMods(&cfg.rootConfig.KeyMaps.ContextKeyMap, modMap)
		arg := enterMod.Arg
		subtitle := enterMod.Subtitle
		vals := enterMod.Variables
		cfg.rootConfig.Awf().Append(
			alfred.NewItem().
				SetTitle(title).
				SetSubtitle(subtitle).
				SetArg(arg).
				SetVariables(vals).
				SetMods(mods),
		)
	}

	return nil
}

func (cfg *Config) GetQuery() string {
	return cfg.fs.Arg(0)
}
