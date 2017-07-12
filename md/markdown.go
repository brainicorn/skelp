package main

import (
	"os"

	"github.com/brainicorn/skelp/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {

	dir := "./"

	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	root := cmd.NewSkelpCommand()
	doc.GenMarkdownTree(root, dir)
}
