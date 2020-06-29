package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/konoui/alfred-k8s/pkg/kubectl"
	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
)

const (
	allNamespacesFlag = "all"
	useFalg           = "use"
	deleteFlag        = "delete"
	namespaceFlag     = "namespace"
)

type middlewareFunc func(*cobra.Command, []string) error

func outputMiddleware(f middlewareFunc) middlewareFunc {
	return func(cmd *cobra.Command, args []string) (ret error) {
		// Note currently we assume for `obj` command
		query := getQuery(args, 1)
		_ = f(cmd, args)
		awf.Filter(query).Output()
		return
	}
}

func cacheOutputMiddleware(f middlewareFunc) middlewareFunc {
	// Note always return nil
	return func(cmd *cobra.Command, args []string) (ret error) {
		nonNs := getBoolFlag(cmd, allNamespacesFlag)
		key := getCacheKey(cmd.Name(), !nonNs)
		query := getQuery(args, 0)
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
		defer func() { _ = deleteAllCaches() }()
		return f(cmd, args)
	}
}

// deleteCache delete all resources for current namespace/context resources.
// 1. list pods in current ns
// 2. switch ns
// 3. list pods in switched current ns
func deleteAllCaches() error {
	files, err := ioutil.ReadDir(cacheDir)
	if err != nil {
		return fmt.Errorf("invalid cache directory %s", cacheDir)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if !strings.HasSuffix(f.Name(), cacheSuffix) {
			continue
		}

		path := filepath.Join(cacheDir, f.Name())
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("failed to delete %s", path)
		}
	}
	return nil
}

func getBoolFlag(cmd *cobra.Command, name string) bool {
	v, err := cmd.PersistentFlags().GetBool(name)
	if err != nil {
		return false
	}
	return v
}

func getStringFlag(cmd *cobra.Command, name string) string {
	v, err := cmd.PersistentFlags().GetString(name)
	if err != nil {
		return ""
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

func addNamespaceFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(namespaceFlag, "", "namespace")
}

func getSternMod(i interface{}) *alfred.Mod {
	name, ns := kubectl.GetNameNamespace(i)
	arg := fmt.Sprintf("stern %s", name)
	if ns != "" {
		arg = fmt.Sprintf("%s --namespace %s", arg, ns)
	}
	return &alfred.Mod{
		Subtitle: "copy simple stern command",
		Arg:      arg,
	}
}

func getDeleteMod(cmdName string, i interface{}) *alfred.Mod {
	name, ns := kubectl.GetNameNamespace(i)
	arg := fmt.Sprintf("%s %s --%s", cmdName, name, deleteFlag)
	if ns != "" {
		arg = fmt.Sprintf("%s --%s %s", arg, namespaceFlag, ns)
	}
	return &alfred.Mod{
		Subtitle: "delete it",
		Arg:      arg,
		Variables: map[string]string{
			nextActionKey: nextActionShell,
		},
	}
}

func getUseMod(cmdName string, i interface{}) *alfred.Mod {
	name, ns := kubectl.GetNameNamespace(i)
	arg := fmt.Sprintf("%s %s --%s", cmdName, name, useFalg)
	if ns != "" {
		arg = fmt.Sprintf("%s --%s %s", arg, namespaceFlag, ns)
	}
	return &alfred.Mod{
		Subtitle: "switch to it",
		Arg:      arg,
		Variables: map[string]string{
			nextActionKey: nextActionShell,
		},
	}
}

func getCopyPortForwardMod(res string, i interface{}) *alfred.Mod {
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
		return &alfred.Mod{
			Subtitle: "the resource has no ports",
		}
	}

	arg := fmt.Sprintf("kubectl port-forward %s/%s %s", res, name, strings.Join(ports, " "))
	if ns != "" {
		arg = fmt.Sprintf("%s --namespace %s", arg, ns)
	}

	return &alfred.Mod{
		Subtitle: "copy " + arg,
		Arg:      arg,
	}
}

func getExecPortForwardMod(res string, i interface{}) *alfred.Mod {
	name, ns := kubectl.GetNameNamespace(i)
	arg := fmt.Sprintf("port-forward %s/%s --%s", res, name, useFalg)
	if ns != "" {
		arg = fmt.Sprintf("%s --%s %s", arg, namespaceFlag, ns)
	}

	return &alfred.Mod{
		Subtitle: "port-forward in background",
		Arg:      arg,
		Variables: map[string]string{
			nextActionKey: nextActionJob,
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

// getCacheKey return cache key.
// 1st arg is resource name. 2nd arg means whether namespaced resource or not
func getCacheKey(name string, namespaced bool) string {
	if namespaced {
		return fmt.Sprintf("%s-%s", cacheNamespacedPrefix, name)
	}
	return fmt.Sprintf("%s-%s", cacheNonNamespacedPrefix, name)
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
