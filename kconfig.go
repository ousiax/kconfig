package main

import (
	"os"

	"github.com/qqbuby/kconfig/cmd"
)

func main() {
	root := cmd.NewCmdKonfig()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
