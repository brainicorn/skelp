package main

import (
	"os"

	"github.com/brainicorn/skelp/cmd"
)

func main() {
	code := cmd.Execute(os.Args[1:], nil)

	if code != 0 {
		panic(code)
	}
}
