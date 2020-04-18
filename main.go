package main

import (
	"os"

	"github.com/jimmino/countbeat/cmd"

	_ "github.com/jimmino/countbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
