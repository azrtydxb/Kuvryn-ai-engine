package manifestcheck

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/azrtydxb/kuvryn-ai-engine/internal/engineimage"
)

type Result struct {
	Manifests []string
	Errors    []error
}

func CheckRoot(root string) Result {
	if root == "" {
		root = "."
	}
	var result Result
	_ = filepath.WalkDir(filepath.Join(root, "engines"), func(path string, d os.DirEntry, err error) error {
		if err != nil {
			result.Errors = append(result.Errors, err)
			return nil
		}
		if d.IsDir() || filepath.Base(path) != "image.yaml" {
			return nil
		}
		result.Manifests = append(result.Manifests, filepath.ToSlash(path))
		if err := checkManifest(root, path); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("%s: %w", rel(root, path), err))
		}
		return nil
	})
	sort.Strings(result.Manifests)
	if len(result.Manifests) == 0 {
		result.Errors = append(result.Errors, fmt.Errorf("no engine manifests found under %s", filepath.Join(root, "engines")))
	}
	return result
}

func checkManifest(root, path string) error {
	name, err := manifestName(path)
	if err != nil {
		return err
	}
	if disallowedFlavor(name) {
		return fmt.Errorf("disallowed aggregate flavor %q", name)
	}
	manifest, err := engineimage.LoadManifest(path)
	if err != nil {
		return err
	}
	if disallowedFlavor(manifest.Name) {
		return fmt.Errorf("disallowed aggregate flavor %q", manifest.Name)
	}
	if manifest.Publish.Certification == "hardware-untested" && (manifest.Vendor != "amd" && manifest.Vendor != "intel") {
		return fmt.Errorf("hardware-untested certification is currently allowed only for amd/intel, got %q", manifest.Vendor)
	}
	if manifest.Publish.Certification == "hardware-certified" && (manifest.Vendor == "amd" || manifest.Vendor == "intel") {
		return fmt.Errorf("%s must stay hardware-untested until matching CT runners exist", manifest.Name)
	}
	for _, upstream := range manifest.Upstreams {
		if upstream.Version == "bootstrap" || upstream.Ref == "bootstrap" || upstream.Digest == "sha256:bootstrap" {
			return fmt.Errorf("upstream %q uses placeholder pin", upstream.Name)
		}
	}
	dockerfile := manifest.Build.Dockerfile
	if !filepath.IsAbs(dockerfile) {
		dockerfile = filepath.Join(root, dockerfile)
	}
	if _, err := os.Stat(dockerfile); err != nil {
		return fmt.Errorf("dockerfile %q: %w", manifest.Build.Dockerfile, err)
	}
	return nil
}

func manifestName(path string) (string, error) {
	var header struct {
		Name string `yaml:"name"`
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	if err := yaml.Unmarshal(b, &header); err != nil {
		return "", err
	}
	return header.Name, nil
}

func disallowedFlavor(name string) bool {
	return strings.Contains(name, "all-engines") || strings.Contains(name, "all-gpu-vendors") || strings.HasSuffix(name, "-ray")
}

func rel(root, path string) string {
	r, err := filepath.Rel(root, path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(r)
}
