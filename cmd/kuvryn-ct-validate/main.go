package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/azrtydxb/kuvryn-ai-engine/internal/ct"
	"github.com/azrtydxb/kuvryn-ai-engine/internal/engineimage"
)

func main() {
	manifestPath := flag.String("manifest", "", "path to image.yaml")
	reportPath := flag.String("report", "", "path to CT report JSON")
	flag.Parse()
	if *manifestPath == "" || *reportPath == "" {
		fmt.Fprintln(os.Stderr, "usage: kuvryn-ct-validate --manifest engines/<engine>/<vendor>/image.yaml --report ct-reports/<flavor>.json")
		os.Exit(2)
	}
	manifest, err := engineimage.LoadManifest(*manifestPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	report, err := ct.LoadEvidenceReport(*reportPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := ct.ValidatePublishEvidence(manifest, report); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("CT evidence accepted for %s (%s)\n", manifest.Name, manifest.Publish.Certification)
}
