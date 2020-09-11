package utils

import (
	"fmt"

	"github.com/konoui/alfred-k8s/pkg/kubectl"
	"github.com/konoui/go-alfred"
)

const (
	AllNamespacesFlag = "a"
	UseFlag           = "use"
	DeleteFlag        = "delete"
	NamespaceFlag     = "namespace"
)

// decide next action for workflow filter
const (
	NextActionKey   = "nextAction"
	NextActionCmd   = "cmd"
	NextActionShell = "shell"
	NextActionJob   = "job"
)

func GetCacheKey(name string, namespaced bool) string {
	if namespaced {
		return fmt.Sprintf("namespace-%s", name)
	}
	return fmt.Sprintf("non-namespace-%s", name)
}

func GetSternMod(i interface{}) *alfred.Mod {
	name, ns := kubectl.GetNameNamespace(i)
	arg := fmt.Sprintf("stern %s", name)
	if ns != "" {
		arg = fmt.Sprintf("%s --namespace %s", arg, ns)
	}

	return alfred.NewMod().
		SetSubtitle("copy simple stern command").
		SetArg(arg)
}

func GetCopyMod(subtitle, arg string) *alfred.Mod {
	return alfred.NewMod().
		SetSubtitle(subtitle).
		SetArg(arg)
}

func GetDeleteMod(cmdName string, i interface{}) *alfred.Mod {
	name, ns := kubectl.GetNameNamespace(i)
	arg := fmt.Sprintf("--%s %s", DeleteFlag, name)
	if ns != "" {
		arg = fmt.Sprintf("--%s %s %s", NamespaceFlag, ns, arg)
	}
	cmd := cmdName + " " + arg

	return alfred.NewMod().
		SetSubtitle("delete it").
		SetArg(cmd).
		SetVariable(NextActionKey, NextActionShell)

}

func GetUseMod(cmdName string, i interface{}) *alfred.Mod {
	name, ns := kubectl.GetNameNamespace(i)
	arg := fmt.Sprintf("--%s %s", UseFlag, name)
	if ns != "" {
		arg = fmt.Sprintf("--%s %s %s", NamespaceFlag, ns, arg)
	}
	cmd := cmdName + " " + arg

	return alfred.NewMod().
		SetSubtitle("switch to it").
		SetArg(cmd).
		SetVariable(NextActionKey, NextActionShell)
}

func GetNamespacedResourceTitle(i interface{}) string {
	name, ns := kubectl.GetNameNamespace(i)
	if ns == "" {
		return name
	}
	return fmt.Sprintf("[%s] %s", ns, name)
}

func GetQuery(args []string, idx int) string {
	if len(args) > idx {
		return args[idx]
	}
	return ""
}
