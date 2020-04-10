package cmd

import (
	"fmt"
	"os"
	"reflect"

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

func getNamespaceResourceTitle(i interface{}) string {
	rv := reflect.Indirect(reflect.ValueOf(i))
	rt := rv.Type()
	if _, ok := rt.FieldByName("Name"); !ok {
		// Note unexpected case
		return "UnknownName"
	}

	name := rv.FieldByName("Name").String()
	if _, ok := rt.FieldByName("Namespace"); !ok {
		return name
	}

	ns := rv.FieldByName("Namespace").String()
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
