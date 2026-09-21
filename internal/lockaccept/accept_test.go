package lockaccept

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/azrtydxb/kuvryn-ai-engine/internal/engineimage"
	"gopkg.in/yaml.v3"
)

func TestAcceptWritesLockAfterEvidencePasses(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "VERSION"), []byte("0.1.0\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, "engines", "vllm", "amd")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(dir, "image.yaml")
	writeYAML(t, manifestPath, engineimage.ImageManifest{
		Schema:       engineimage.ManifestSchema,
		Name:         "vllm-amd",
		Engine:       "vllm",
		Vendor:       "amd",
		Image:        engineimage.Image{Repository: engineimage.RepositoryPrefix + "vllm-amd", Platforms: []string{"linux/amd64"}, Base: engineimage.Base{Image: "ubuntu", Tag: "24.04"}},
		Build:        engineimage.Build{Context: ".", Dockerfile: "engines/vllm/amd/Dockerfile"},
		Capabilities: engineimage.Capabilities{Accelerators: []string{"amd"}, Collective: "rccl", Distributed: true},
		Upstreams:    []engineimage.Upstream{{Name: "engine", Type: "static-version", Version: "v1"}},
		Publish:      engineimage.Publish{Certification: "hardware-untested"},
		Tests: []engineimage.TestDeclaration{
			{Name: "static", Type: "static", Required: true},
			{Name: "vendor-device", Type: "vendor-device", Required: false},
		},
	})
	reportPath := filepath.Join(root, "report.json")
	if err := os.WriteFile(reportPath, []byte(`{"schema":"kuvryn.container-test-report/v1","flavor":"vllm-amd","image":"example","results":[{"name":"static","type":"static","status":"passed","required":true}]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(root, "accepted.lock")
	lock, written, err := Accept(context.Background(), Options{Root: root, ManifestPath: manifestPath, ReportPath: reportPath, ImageDigest: "sha256:built", OutputPath: out})
	if err != nil {
		t.Fatal(err)
	}
	if written != out || lock.ImageDigest != "sha256:built" || len(lock.ResolvedInputs) != 1 {
		t.Fatalf("lock=%#v written=%s", lock, written)
	}
	if len(lock.PublishedTags) != 4 || lock.PublishedTags[1] != "v0.1.0_amd_vllm_1" || lock.PublishedTags[3] != "hardware-untested" {
		t.Fatalf("lock=%#v written=%s", lock, written)
	}
	loaded, err := engineimage.LoadLock(out)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Flavor != "vllm-amd" || loaded.PublishedTags[1] != "v0.1.0_amd_vllm_1" || loaded.TestEvidence[0].Status != "passed" {
		t.Fatalf("loaded=%#v", loaded)
	}
}

func TestAcceptRejectsBadEvidence(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "engines", "vllm", "nvidia")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(dir, "image.yaml")
	writeYAML(t, manifestPath, engineimage.ImageManifest{
		Schema:       engineimage.ManifestSchema,
		Name:         "vllm-nvidia",
		Engine:       "vllm",
		Vendor:       "nvidia",
		Image:        engineimage.Image{Repository: engineimage.RepositoryPrefix + "vllm-nvidia", Platforms: []string{"linux/amd64"}, Base: engineimage.Base{Image: "ubuntu", Tag: "24.04"}},
		Build:        engineimage.Build{Context: ".", Dockerfile: "engines/vllm/nvidia/Dockerfile"},
		Capabilities: engineimage.Capabilities{Accelerators: []string{"nvidia"}, Collective: "nccl", Distributed: true},
		Upstreams:    []engineimage.Upstream{{Name: "engine", Type: "static-version", Version: "v1"}},
		Publish:      engineimage.Publish{Certification: "hardware-certified"},
		Tests: []engineimage.TestDeclaration{
			{Name: "static", Type: "static", Required: true},
			{Name: "vendor-device", Type: "vendor-device", Required: true},
		},
	})
	reportPath := filepath.Join(root, "report.json")
	if err := os.WriteFile(reportPath, []byte(`{"schema":"kuvryn.container-test-report/v1","flavor":"vllm-nvidia","image":"example","results":[{"name":"static","type":"static","status":"passed","required":true}]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := Accept(context.Background(), Options{Root: root, ManifestPath: manifestPath, ReportPath: reportPath, ImageDigest: "sha256:built"}); err == nil {
		t.Fatal("expected missing vendor CT error")
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
