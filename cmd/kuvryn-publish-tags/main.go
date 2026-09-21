package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/azrtydxb/kuvryn-ai-engine/internal/engineimage"
	"github.com/azrtydxb/kuvryn-ai-engine/internal/publishtags"
)

func main() {
	manifestPath := flag.String("manifest", "", "engine image manifest path")
	digest := flag.String("digest", "", "published image digest sha256:<hex>")
	version := flag.String("version", "", "Kuvryn AI Engine version; defaults to --version-file")
	versionFile := flag.String("version-file", "VERSION", "file containing Kuvryn AI Engine version")
	format := flag.String("format", "json", "output format: json or args")
	flag.Parse()
	if *manifestPath == "" || *digest == "" {
		fmt.Fprintln(os.Stderr, "--manifest and --digest are required")
		os.Exit(2)
	}
	appVersion := strings.TrimSpace(*version)
	if appVersion == "" {
		b, err := os.ReadFile(*versionFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read version file: %v\n", err)
			os.Exit(2)
		}
		appVersion = strings.TrimSpace(string(b))
	}
	manifest, err := engineimage.LoadManifest(*manifestPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	plan, err := publishtags.Build(manifest, *digest, appVersion)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	switch *format {
	case "json":
		b, err := json.MarshalIndent(plan, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(string(b))
	case "args":
		for _, tag := range plan.Tags {
			fmt.Printf("-t\n%s:%s\n", plan.Repository, tag)
		}
	default:
		fmt.Fprintf(os.Stderr, "unsupported format %q\n", *format)
		os.Exit(2)
	}
}
