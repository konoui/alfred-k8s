package main

import (
	"github.com/konoui/alfred-k8s/cmd"
)

func main() {
	c := cmd.NewCmd()
	cmd.Execute(c)
}
