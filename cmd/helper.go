package cmd

import (
	"fmt"
	"os"

	"github.com/konoui/alfred-k8s/pkg/kubectl"
	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

const (
	allNamespacesFlag = "all"
	useFalg           = "use"
	deleteFlag        = "delete"
)

type middlewareFunc func(*cobra.Command, []string) error

func outputMiddleware(f middlewareFunc) middlewareFunc {
	// Note always return nil
	return func(cmd *cobra.Command, args []string) (ret error) {
		all := getBoolFlag(cmd, allNamespacesFlag)
		query := getQuery(args, 0)
		key := getCacheKey(cmd.Name(), all)
		if awf.Cache(key).MaxAge(cacheTime).LoadItems().Err() == nil {
			awf.Filter(query).Output()
			return
		}
		if err := f(cmd, args); err != nil {
			fatal(err)
			return
		}
		awf.Cache(key).StoreItems().Workflow().Filter(query).Output()
		return
	}
}

func shellOutputMiddleware(f middlewareFunc) middlewareFunc {
	// Note always return nil
	return func(cmd *cobra.Command, args []string) (ret error) {
		if err := f(cmd, args); err != nil {
			fmt.Fprintf(outStream, "Failed due to %s\n", err)
			return
		}
		fmt.Fprintf(outStream, "Success!!\n")
		return
	}
}

func clearCacheMiddleware(f middlewareFunc) middlewareFunc {
	return func(cmd *cobra.Command, args []string) error {
		all := getBoolFlag(cmd, allNamespacesFlag)
		key := getCacheKey(cmd.Name(), all)
		defer func() { awf.Cache(key).Delete() }()
		return f(cmd, args)
	}
}

func getBoolFlag(cmd *cobra.Command, name string) bool {
	v, err := cmd.PersistentFlags().GetBool(name)
	if err != nil {
		return false
	}
	return v
}

func addAllNamespacesFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP(allNamespacesFlag, "a", false, "in all namespaces")
}

func addUseFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool(useFalg, false, "switch to it")
}

func addDeleteFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool(deleteFlag, false, "delete the resource")
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
		Arg:      fmt.Sprintf("%s %s %s --%s", rs, name, ns, deleteFlag),
		Variables: map[string]string{
			nextActionKey: nextActionShell,
		},
	}
}

func getUseMod(rs string, i interface{}) alfred.Mod {
	name, _ := kubectl.GetNameNamespace(i)
	return alfred.Mod{
		Subtitle: "switch to it",
		Arg:      fmt.Sprintf("%s %s --%s", rs, name, useFalg),
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
		fatal(err)
		os.Exit(255)
	}
}

func fatal(err error) {
	awf.Fatal("Fatal error occurs", err.Error())
}
