package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/azrtydxb/kuvryn-ai-engine/internal/lockaccept"
)

func main() {
	root := flag.String("root", ".", "repository root for resolving local upstream inputs")
	manifest := flag.String("manifest", "", "path to image.yaml")
	report := flag.String("report", "", "path to CT report JSON")
	digest := flag.String("digest", "", "published image digest; defaults to report digest")
	output := flag.String("output", "", "lock output path; defaults to sibling upstream.lock")
	version := flag.String("version", "", "Kuvryn AI Engine version; defaults to --version-file")
	versionFile := flag.String("version-file", "", "file containing Kuvryn AI Engine version; defaults to <root>/VERSION")
	flag.Parse()
	if *manifest == "" || *report == "" {
		fmt.Fprintln(os.Stderr, "usage: kuvryn-lock-accept --manifest engines/<engine>/<vendor>/image.yaml --report ct-reports/<flavor>.json --digest sha256:...")
		os.Exit(2)
	}
	lock, out, err := lockaccept.Accept(context.Background(), lockaccept.Options{Root: *root, ManifestPath: *manifest, ReportPath: *report, ImageDigest: *digest, OutputPath: *output, Version: *version, VersionFile: *versionFile})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("accepted %s lock at %s\n", lock.Flavor, out)
}
