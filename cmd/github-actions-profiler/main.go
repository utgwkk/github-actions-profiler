package main

import (
	"context"
	"os"

	"github.com/utgwkk/github-actions-profiler"
)

func main() {
	ctx := context.Background()
	ghaprofiler.StartCLI(ctx, os.Args[1:])
}
