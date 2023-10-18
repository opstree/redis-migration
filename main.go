package main

import (
	"fmt"
	"github.com/opstree/redis-migration/cmd"
)

var (
	GitCommit string
	Version   string
)

func main() {
	version := fmt.Sprintf("%s: %s", Version, GitCommit)
	cmd.Execute(version)
}
