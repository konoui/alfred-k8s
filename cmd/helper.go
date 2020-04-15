package cmd

import (
	"fmt"
	"os"

	"github.com/konoui/alfred-k8s/pkg/kubectl"
	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

func addAllNamespaceFlag(cmd *cobra.Command, all *bool) {
	cmd.PersistentFlags().BoolVarP(all, "all", "a", false, "in all namespaces")
}

func addUseFlag(cmd *cobra.Command, use *bool) {
	cmd.PersistentFlags().BoolVarP(use, "use", "u", false, "switch to it")
}

func addDeleteFlag(cmd *cobra.Command, del *bool) {
	cmd.PersistentFlags().BoolVarP(del, "delete", "d", false, "delete the resource")
}

func getSternMod(i interface{}) alfred.Mod {
	name, ns := kubectl.GetNameNamespace(i)
	arg := fmt.Sprintf("stern %s", name)
	if ns != "" {
		arg = fmt.Sprintf("%s --namespace %s", arg, ns)
	}

	return alfred.Mod{
		Subtitle: "copy simple stern command",
		Arg:      arg,
	}
}

func getNamespaceResourceTitle(i interface{}) string {
	name, ns := kubectl.GetNameNamespace(i)
	if ns == "" {
		return name
	}

	return fmt.Sprintf("[%s] %s", ns, name)
}

func getQuery(args []string, idx int) string {
	if len(args) > idx {
		return args[idx]
	}
	return ""
}

func exitWith(err error) {
	if err != nil {
		awf.Fatal("fatal error occurs", err.Error())
		os.Exit(255)
	}
}
