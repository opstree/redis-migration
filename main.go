package main

import (
	"fmt"
	"redis-migrator/cmd"
)

var (
	GitCommit string
	Version   string
)

func main() {
	version := fmt.Sprintf("%s: %s", Version, GitCommit)
	cmd.Execute(version)
}
