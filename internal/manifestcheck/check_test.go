package manifestcheck

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/azrtydxb/kuvryn-ai-engine/internal/engineimage"
	"gopkg.in/yaml.v3"
)

func TestCheckRootAcceptsValidManifest(t *testing.T) {
	root := t.TempDir()
	writeFlavor(t, root, "llama-cpp-cpu", "llama-cpp", "cpu", "cpu-only")
	result := CheckRoot(root)
	if len(result.Errors) != 0 {
		t.Fatalf("errors = %v", result.Errors)
	}
	if len(result.Manifests) != 1 {
		t.Fatalf("manifests = %v", result.Manifests)
	}
}

func TestCheckRootRejectsDisallowedFlavor(t *testing.T) {
	root := t.TempDir()
	writeFlavor(t, root, "all-engines-nvidia", "vllm", "nvidia", "hardware-certified")
	result := CheckRoot(root)
	if len(result.Errors) == 0 || !strings.Contains(result.Errors[0].Error(), "disallowed aggregate flavor") {
		t.Fatalf("errors = %v", result.Errors)
	}
}

func TestCheckRootRejectsAMDCertified(t *testing.T) {
	root := t.TempDir()
	writeFlavor(t, root, "vllm-amd", "vllm", "amd", "hardware-certified")
	result := CheckRoot(root)
	if len(result.Errors) == 0 || !strings.Contains(result.Errors[0].Error(), "must be hardware-untested") {
		t.Fatalf("errors = %v", result.Errors)
	}
}

func TestCheckRootRejectsPlaceholderPins(t *testing.T) {
	root := t.TempDir()
	writeFlavor(t, root, "llama-cpp-cpu", "llama-cpp", "cpu", "cpu-only", func(m *engineimage.ImageManifest) {
		m.Upstreams[0].Version = "bootstrap"
	})
	result := CheckRoot(root)
	if len(result.Errors) == 0 || !strings.Contains(result.Errors[0].Error(), "placeholder pin") {
		t.Fatalf("errors = %v", result.Errors)
	}
}

func writeFlavor(t *testing.T, root, name, engine, vendor, certification string, mutate ...func(*engineimage.ImageManifest)) {
	t.Helper()
	dir := filepath.Join(root, "engines", engine, vendor)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	dockerfile := filepath.Join("engines", engine, vendor, "Dockerfile")
	if err := os.WriteFile(filepath.Join(root, dockerfile), []byte("FROM scratch\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	manifest := engineimage.ImageManifest{
		Schema:       engineimage.ManifestSchema,
		Name:         name,
		Engine:       engine,
		Vendor:       vendor,
		Image:        engineimage.Image{Repository: engineimage.RepositoryPrefix + name, Platforms: []string{"linux/amd64"}, Base: engineimage.Base{Image: "ubuntu", Tag: "24.04"}},
		Build:        engineimage.Build{Context: ".", Dockerfile: dockerfile},
		Capabilities: engineimage.Capabilities{Accelerators: []string{vendor}, Collective: "none"},
		Upstreams:    []engineimage.Upstream{{Name: "engine", Type: "static-version", Version: "test"}},
		Publish:      engineimage.Publish{Certification: certification},
		Tests:        []engineimage.TestDeclaration{{Name: "static", Type: "static", Required: true}},
	}
	if vendor == "cpu" {
		manifest.Capabilities.Accelerators = []string{"cpu"}
	}
	for _, fn := range mutate {
		fn(&manifest)
	}
	b, err := yaml.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "image.yaml"), b, 0o600); err != nil {
		t.Fatal(err)
	}
}
