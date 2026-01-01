package main

import (
	"os"

	"github.com/gizzahub/gzh-cli-gitflow/cmd/gz-flow/cmd"
)

var version = "dev"

func main() {
	cmd.SetVersion(version)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
