package main

import (
	"github.com/konoui/alfred-k8s/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	cmd.Execute(rootCmd)
}
