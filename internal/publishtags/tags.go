package publishtags

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/azrtydxb/kuvryn-ai-engine/internal/engineimage"
)

var tagUnsafe = regexp.MustCompile(`[^A-Za-z0-9_.-]+`)

type Plan struct {
	Repository string   `json:"repository"`
	Tags       []string `json:"tags"`
}

func Build(manifest engineimage.ImageManifest, digest, appVersion string) (Plan, error) {
	if digest == "" || !strings.HasPrefix(digest, "sha256:") {
		return Plan{}, fmt.Errorf("digest must be sha256:<hex>")
	}
	engineVersion, err := engineVersion(manifest)
	if err != nil {
		return Plan{}, err
	}
	appVersion = normalizeAppVersion(appVersion)
	engine := sanitize(manifest.Engine)
	vendor := sanitize(manifest.Vendor)
	engineVersion = sanitize(strings.TrimPrefix(engineVersion, "v"))
	if appVersion == "v" || engine == "" || vendor == "" || engineVersion == "" {
		return Plan{}, fmt.Errorf("cannot build publish tags from app=%q engine=%q vendor=%q engineVersion=%q", appVersion, manifest.Engine, manifest.Vendor, engineVersion)
	}
	tags := []string{
		digestTag(digest),
		fmt.Sprintf("%s_%s_%s_%s", appVersion, vendor, engine, engineVersion),
		fmt.Sprintf("%s_%s_%s_nightly", appVersion, vendor, engine),
	}
	switch manifest.Publish.Certification {
	case "hardware-certified":
		tags = append(tags, "latest", "hardware-certified")
	case "cpu-only":
		tags = append(tags, "latest", "cpu-only")
	case "hardware-untested":
		tags = append(tags, "hardware-untested")
	default:
		return Plan{}, fmt.Errorf("unsupported certification %q", manifest.Publish.Certification)
	}
	return Plan{Repository: manifest.Image.Repository, Tags: unique(tags)}, nil
}

func engineVersion(manifest engineimage.ImageManifest) (string, error) {
	for _, upstream := range manifest.Upstreams {
		if upstream.Name != "engine" {
			continue
		}
		if upstream.Version != "" {
			return upstream.Version, nil
		}
		if upstream.Ref != "" {
			return upstream.Ref, nil
		}
		if upstream.Digest != "" {
			return upstream.Digest, nil
		}
	}
	return "", fmt.Errorf("manifest %q does not declare an engine upstream version/ref", manifest.Name)
}

func normalizeAppVersion(version string) string {
	version = strings.TrimSpace(version)
	version = strings.TrimPrefix(version, "v")
	return "v" + sanitize(version)
}

func digestTag(digest string) string { return strings.TrimPrefix(digest, "sha256:") }

func sanitize(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, "_-")
	value = tagUnsafe.ReplaceAllString(value, "_")
	return value
}

func unique(values []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	return out
}
