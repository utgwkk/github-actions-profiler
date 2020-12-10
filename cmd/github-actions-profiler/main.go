package main

import (
	"context"
	"os"

	ghaprofiler "github.com/utgwkk/github-actions-profiler"
)

func main() {
	ctx := context.Background()
	cli := ghaprofiler.NewCLI()
	cli.Start(ctx, os.Args[1:])
}
