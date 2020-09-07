package rootcmd

import (
	"context"
	"errors"
	"io"

	"github.com/konoui/alfred-k8s/pkg/kubectl"
	"github.com/konoui/go-alfred"
	"github.com/peterbourgon/ff/v3/ffcli"
)

type Config struct {
	outStream io.Writer
	errStream io.Writer
	k         *kubectl.Kubectl
	awf       *alfred.Workflow
}

func NewConfig(out, err io.Writer, k *kubectl.Kubectl, awf *alfred.Workflow) *Config {
	return &Config{
		outStream: out,
		errStream: err,
		k:         k,
		awf:       awf,
	}
}

func (cfg *Config) Stdout() io.Writer {
	return cfg.outStream
}

func (cfg *Config) Stderr() io.Writer {
	return cfg.errStream
}

func (cfg *Config) Kubeclt() *kubectl.Kubectl {
	return cfg.k
}

func (cfg *Config) Awf() *alfred.Workflow {
	return cfg.awf
}

// New create a new cmd for root
func New() *ffcli.Command {
	rootCmd := &ffcli.Command{
		Name: "list available commands",
		Exec: func(ctx context.Context, args []string) error {
			return errors.New("not implemented")
		},
	}

	return rootCmd
}
