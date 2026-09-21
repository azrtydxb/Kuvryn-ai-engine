package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/azrtydxb/kuvryn-ai-engine/internal/detect"
)

func main() {
	root := flag.String("root", ".", "repository root")
	output := flag.String("output", "", "write update plan JSON to this path instead of stdout")
	flag.Parse()

	plan, err := detect.Detect(context.Background(), detect.Options{Root: *root})
	b, marshalErr := json.MarshalIndent(plan, "", "  ")
	if marshalErr != nil {
		fmt.Fprintln(os.Stderr, marshalErr)
		os.Exit(1)
	}
	b = append(b, '\n')
	if *output == "" {
		_, _ = os.Stdout.Write(b)
	} else if writeErr := os.WriteFile(*output, b, 0o644); writeErr != nil {
		fmt.Fprintln(os.Stderr, writeErr)
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
