package cmd

import (
	"fmt"
	"os"

	"github.com/konoui/alfred-k8s/pkg/kubectl"
	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

const (
	allNamespaceFlag = "--all"
	useFalg          = "--use"
	deleteFlag       = "--delete"
)

func addAllNamespaceFlag(cmd *cobra.Command, all *bool) {
	cmd.PersistentFlags().BoolVarP(all, allNamespaceFlag[2:], "a", false, "in all namespaces")
}

func addUseFlag(cmd *cobra.Command, use *bool) {
	cmd.PersistentFlags().BoolVarP(use, useFalg[2:], "u", false, "switch to it")
}

func addDeleteFlag(cmd *cobra.Command, del *bool) {
	cmd.PersistentFlags().BoolVarP(del, deleteFlag[2:], "d", false, "delete the resource")
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

func getDeleteMod(rs string, i interface{}) alfred.Mod {
	name, ns := kubectl.GetNameNamespace(i)
	return alfred.Mod{
		Subtitle: "delete it",
		Arg:      fmt.Sprintf("%s %s %s %s", rs, name, ns, deleteFlag),
		Variables: map[string]string{
			nextActionKey: nextActionShell,
		},
	}
}

func getUseMod(rs string, i interface{}) alfred.Mod {
	name, _ := kubectl.GetNameNamespace(i)
	return alfred.Mod{
		Subtitle: "switch to it",
		Arg:      fmt.Sprintf("%s %s %s", rs, name, useFalg),
		Variables: map[string]string{
			nextActionKey: nextActionShell,
		},
	}
}

func getNamespaceResourceTitle(i interface{}) string {
	name, ns := kubectl.GetNameNamespace(i)
	if ns == "" {
		return name
	}
	return fmt.Sprintf("[%s] %s", ns, name)
}

func getCacheKey(name string, ns bool) string {
	key := name
	if ns {
		key = fmt.Sprintf("%s-in-all-ns", name)
	}
	return key
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
