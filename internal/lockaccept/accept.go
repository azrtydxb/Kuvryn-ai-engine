package lockaccept

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/azrtydxb/kuvryn-ai-engine/internal/ct"
	"github.com/azrtydxb/kuvryn-ai-engine/internal/detect"
	"github.com/azrtydxb/kuvryn-ai-engine/internal/engineimage"
	"github.com/azrtydxb/kuvryn-ai-engine/internal/publishtags"
)

type Options struct {
	Root         string
	ManifestPath string
	ReportPath   string
	ImageDigest  string
	OutputPath   string
	Version      string
	VersionFile  string
}

func Accept(ctx context.Context, opts Options) (engineimage.UpstreamLock, string, error) {
	root := opts.Root
	if root == "" {
		root = "."
	}
	manifest, err := engineimage.LoadManifest(opts.ManifestPath)
	if err != nil {
		return engineimage.UpstreamLock{}, "", err
	}
	report, err := ct.LoadEvidenceReport(opts.ReportPath)
	if err != nil {
		return engineimage.UpstreamLock{}, "", err
	}
	if err := ct.ValidatePublishEvidence(manifest, report); err != nil {
		return engineimage.UpstreamLock{}, "", err
	}
	resolved := make([]engineimage.ResolvedInput, 0, len(manifest.Upstreams))
	resolvers := detect.DefaultResolvers(root)
	for _, upstream := range manifest.Upstreams {
		r, err := resolvers.Resolve(ctx, upstream)
		if err != nil {
			return engineimage.UpstreamLock{}, "", err
		}
		resolved = append(resolved, r)
	}
	digest := opts.ImageDigest
	if digest == "" {
		digest = report.Digest
	}
	version, err := appVersion(root, opts)
	if err != nil {
		return engineimage.UpstreamLock{}, "", err
	}
	tagPlan, err := publishtags.Build(manifest, digest, version)
	if err != nil {
		return engineimage.UpstreamLock{}, "", err
	}
	lock := engineimage.UpstreamLock{
		Schema:         engineimage.LockSchema,
		Flavor:         manifest.Name,
		ResolvedInputs: resolved,
		ImageDigest:    digest,
		PublishedTags:  tagPlan.Tags,
		TestEvidence:   evidenceFromReport(report),
	}
	output := opts.OutputPath
	if output == "" {
		output = filepath.Join(filepath.Dir(opts.ManifestPath), engineimage.DefaultLockFilename)
	}
	if err := writeLock(output, lock); err != nil {
		return engineimage.UpstreamLock{}, "", err
	}
	return lock, output, nil
}

func appVersion(root string, opts Options) (string, error) {
	if strings.TrimSpace(opts.Version) != "" {
		return strings.TrimSpace(opts.Version), nil
	}
	versionFile := opts.VersionFile
	if versionFile == "" {
		versionFile = filepath.Join(root, "VERSION")
	}
	b, err := os.ReadFile(versionFile)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

func evidenceFromReport(report ct.EvidenceReport) []engineimage.TestEvidence {
	results := report.AllResults()
	evidence := make([]engineimage.TestEvidence, 0, len(results))
	for _, result := range results {
		evidence = append(evidence, engineimage.TestEvidence{Name: result.Name, Status: result.Status})
	}
	return evidence
}

func writeLock(path string, lock engineimage.UpstreamLock) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := yaml.Marshal(lock)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}
