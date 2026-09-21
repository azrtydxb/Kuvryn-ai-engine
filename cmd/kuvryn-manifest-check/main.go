package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/azrtydxb/kuvryn-ai-engine/internal/manifestcheck"
)

func main() {
	root := flag.String("root", ".", "repository root")
	flag.Parse()
	result := manifestcheck.CheckRoot(*root)
	for _, err := range result.Errors {
		fmt.Fprintln(os.Stderr, err)
	}
	if len(result.Errors) > 0 {
		os.Exit(1)
	}
	fmt.Printf("validated %d engine manifest(s)\n", len(result.Manifests))
}
