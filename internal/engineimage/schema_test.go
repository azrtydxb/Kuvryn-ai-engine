package engineimage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadManifestValid(t *testing.T) {
	path := writeTemp(t, "image.yaml", validManifest("llama-cpp-cpu", "llama-cpp", "cpu"))
	m, err := LoadManifest(path)
	if err != nil {
		t.Fatalf("LoadManifest() error = %v", err)
	}
	if m.Name != "llama-cpp-cpu" {
		t.Fatalf("Name = %q", m.Name)
	}
}

func TestLoadManifestRejectsUnknownField(t *testing.T) {
	path := writeTemp(t, "image.yaml", validManifest("llama-cpp-cpu", "llama-cpp", "cpu")+"unexpected: true\n")
	_, err := LoadManifest(path)
	if err == nil || !strings.Contains(err.Error(), "field unexpected not found") {
		t.Fatalf("LoadManifest() error = %v, want unknown field", err)
	}
}

func TestValidateManifestFailures(t *testing.T) {
	cases := map[string]ImageManifest{
		"bad namespace": manifestStruct("vllm-nvidia", "vllm", "nvidia", func(m *ImageManifest) { m.Image.Repository = "ghcr.io/example/vllm-nvidia" }),
		"bad upstream":  manifestStruct("vllm-nvidia", "vllm", "nvidia", func(m *ImageManifest) { m.Upstreams = []Upstream{{Name: "x", Type: "mystery"}} }),
		"no tests":      manifestStruct("vllm-nvidia", "vllm", "nvidia", func(m *ImageManifest) { m.Tests = nil }),
		"bad combo":     manifestStruct("tensorrt-llm-amd", "tensorrt-llm", "amd", nil),
	}
	for name, m := range cases {
		t.Run(name, func(t *testing.T) {
			if err := ValidateManifest(m); err == nil {
				t.Fatal("ValidateManifest() nil error")
			}
		})
	}
}

func TestLoadLockMissingAndMalformed(t *testing.T) {
	if _, err := LoadLock(filepath.Join(t.TempDir(), "missing.lock")); err == nil {
		t.Fatal("LoadLock(missing) nil error")
	}
	path := writeTemp(t, "upstream.lock", "schema: kuvryn.engine-lock/v1\nflavor: [not-string]\n")
	if _, err := LoadLock(path); err == nil {
		t.Fatal("LoadLock(malformed) nil error")
	}
}

func TestExampleLocksLoad(t *testing.T) {
	for _, path := range []string{
		filepath.Join("..", "..", "examples", "locks", "llama-cpp-cpu.upstream.lock"),
		filepath.Join("..", "..", "examples", "locks", "vllm-nvidia.upstream.lock"),
	} {
		lock, err := LoadLock(path)
		if err != nil {
			t.Fatalf("LoadLock(%s) error = %v", path, err)
		}
		if len(lock.PublishedTags) == 0 {
			t.Fatalf("LoadLock(%s) missing publishedTags", path)
		}
	}
}

func validManifest(name, engine, vendor string) string {
	return `schema: kuvryn.engine-image/v1
name: ` + name + `
engine: ` + engine + `
vendor: ` + vendor + `
image:
  repository: ghcr.io/azrtydxb/kuvryn-ai-engine/` + name + `
  platforms: [linux/amd64]
  base:
    image: ubuntu
    tag: "24.04"
build:
  context: .
  dockerfile: Dockerfile
capabilities:
  accelerators: [` + vendor + `]
  collective: none
  distributed: false
upstreams:
  - name: engine
    type: static-version
    version: test
publish:
  certification: cpu-only
tests:
  - name: static
    type: static
    required: true
`
}

func manifestStruct(name, engine, vendor string, mutate func(*ImageManifest)) ImageManifest {
	acc := vendor
	if vendor == "cpu" {
		acc = "cpu"
	}
	m := ImageManifest{
		Schema:       ManifestSchema,
		Name:         name,
		Engine:       engine,
		Vendor:       vendor,
		Image:        Image{Repository: RepositoryPrefix + name, Platforms: []string{"linux/amd64"}, Base: Base{Image: "ubuntu", Tag: "24.04"}},
		Build:        Build{Context: ".", Dockerfile: "Dockerfile"},
		Capabilities: Capabilities{Accelerators: []string{acc}, Collective: "none"},
		Upstreams:    []Upstream{{Name: "engine", Type: "static-version", Version: "test"}},
		Publish:      Publish{Certification: certificationForVendor(vendor)},
		Tests:        []TestDeclaration{{Name: "static", Type: "static", Required: true}},
	}
	if mutate != nil {
		mutate(&m)
	}
	return m
}

func certificationForVendor(vendor string) string {
	switch vendor {
	case "cpu":
		return "cpu-only"
	case "amd", "intel":
		return "hardware-untested"
	default:
		return "hardware-certified"
	}
}

func writeTemp(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}
