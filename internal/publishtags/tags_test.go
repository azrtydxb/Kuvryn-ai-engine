package publishtags

import (
	"reflect"
	"testing"

	"github.com/azrtydxb/kuvryn-ai-engine/internal/engineimage"
)

func TestBuildTagsHardwareCertified(t *testing.T) {
	plan, err := Build(manifest("vllm", "0.29.0", "hardware-certified"), "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "1.0.2")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "v1.0.2_nvidia_vllm_0.29.0", "v1.0.2_nvidia_vllm_nightly", "latest", "hardware-certified"}
	if !reflect.DeepEqual(plan.Tags, want) {
		t.Fatalf("tags = %#v, want %#v", plan.Tags, want)
	}
}

func TestBuildTagsStripsEngineVersionVPrefix(t *testing.T) {
	plan, err := Build(manifest("llama-cpp", "v0.4.1", "cpu-only"), "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "v1.0.2")
	if err != nil {
		t.Fatal(err)
	}
	if plan.Tags[1] != "v1.0.2_nvidia_llama-cpp_0.4.1" {
		t.Fatalf("version tag = %q", plan.Tags[1])
	}
}

func TestBuildTagsHardwareUntestedDoesNotAddLatest(t *testing.T) {
	plan, err := Build(manifest("vllm", "0.29.0", "hardware-untested"), "sha256:cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc", "1.0.2")
	if err != nil {
		t.Fatal(err)
	}
	for _, tag := range plan.Tags {
		if tag == "latest" {
			t.Fatalf("hardware-untested tags include latest: %#v", plan.Tags)
		}
	}
}

func manifest(engine, version, certification string) engineimage.ImageManifest {
	return engineimage.ImageManifest{
		Name:   engine + "-nvidia",
		Engine: engine,
		Vendor: "nvidia",
		Image:  engineimage.Image{Repository: engineimage.RepositoryPrefix + engine + "-nvidia"},
		Upstreams: []engineimage.Upstream{{
			Name:    "engine",
			Type:    "static-version",
			Version: version,
		}},
		Publish: engineimage.Publish{Certification: certification},
	}
}
