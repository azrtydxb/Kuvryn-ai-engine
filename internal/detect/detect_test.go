package detect

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/azrtydxb/kuvryn-ai-engine/internal/engineimage"
	"gopkg.in/yaml.v3"
)

func TestDetectMarksOnlyChangedFlavor(t *testing.T) {
	root := t.TempDir()
	writeFlavor(t, root, "llama-cpp", "cpu", "same", true)
	writeFlavor(t, root, "vllm", "nvidia", "new", false)

	plan, err := Detect(context.Background(), Options{Root: root})
	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}
	if len(plan.Flavors) != 2 {
		t.Fatalf("len(flavors) = %d", len(plan.Flavors))
	}
	got := map[string]FlavorPlan{}
	for _, f := range plan.Flavors {
		got[f.Name] = f
	}
	if got["llama-cpp-cpu"].Changed {
		t.Fatalf("llama-cpp-cpu changed = true, want false")
	}
	if !got["vllm-nvidia"].Changed || got["vllm-nvidia"].Reason != "upstream-changed" {
		t.Fatalf("vllm-nvidia = changed %v reason %q", got["vllm-nvidia"].Changed, got["vllm-nvidia"].Reason)
	}
	if got["llama-cpp-cpu"].Certification != "cpu-only" {
		t.Fatalf("certification = %q", got["llama-cpp-cpu"].Certification)
	}
}

func TestDetectMissingLockIsChanged(t *testing.T) {
	root := t.TempDir()
	writeFlavor(t, root, "llama-cpp", "cpu", "same", false)
	plan, err := Detect(context.Background(), Options{Root: root})
	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}
	if !plan.Flavors[0].Changed || plan.Flavors[0].Reason != "missing-lock" {
		t.Fatalf("flavor = changed %v reason %q", plan.Flavors[0].Changed, plan.Flavors[0].Reason)
	}
}

func writeFlavor(t *testing.T, root, engine, vendor, manifestVersion string, lockSame bool) {
	t.Helper()
	name := engine + "-" + vendor
	dir := filepath.Join(root, "engines", engine, vendor)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := engineimage.ImageManifest{
		Schema:       engineimage.ManifestSchema,
		Name:         name,
		Engine:       engine,
		Vendor:       vendor,
		Image:        engineimage.Image{Repository: engineimage.RepositoryPrefix + name, Platforms: []string{"linux/amd64"}, Base: engineimage.Base{Image: "ubuntu", Tag: "24.04"}},
		Build:        engineimage.Build{Context: ".", Dockerfile: "Dockerfile"},
		Capabilities: engineimage.Capabilities{Accelerators: []string{vendor}, Collective: "none"},
		Upstreams:    []engineimage.Upstream{{Name: "engine", Type: "static-version", Version: manifestVersion}},
		Publish:      engineimage.Publish{Certification: "hardware-certified"},
		Tests:        []engineimage.TestDeclaration{{Name: "static", Type: "static", Required: true}},
	}
	if vendor == "cpu" {
		manifest.Capabilities.Accelerators = []string{"cpu"}
		manifest.Publish.Certification = "cpu-only"
	}
	if vendor == "amd" || vendor == "intel" {
		manifest.Publish.Certification = "hardware-untested"
	}
	writeYAML(t, filepath.Join(dir, "image.yaml"), manifest)
	if lockSame || manifestVersion == "new" {
		version := manifestVersion
		if manifestVersion == "new" {
			version = "old"
		}
		lock := engineimage.UpstreamLock{Schema: engineimage.LockSchema, Flavor: name, ResolvedInputs: []engineimage.ResolvedInput{{Name: "engine", Type: "static-version", Value: version, Source: "engine"}}}
		writeYAML(t, filepath.Join(dir, "upstream.lock"), lock)
	}
}

func writeYAML(t *testing.T, path string, v any) {
	t.Helper()
	b, err := yaml.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, b, 0o600); err != nil {
		t.Fatal(err)
	}
}
