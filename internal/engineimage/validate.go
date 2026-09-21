package engineimage

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var flavorRe = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*-(?:cpu|nvidia|amd|intel)(?:-[a-z0-9]+)*$`)

var allowedUpstreamTypes = map[string]bool{
	"github-release":   true,
	"github-tag":       true,
	"github-commit":    true,
	"container-digest": true,
	"pypi-version":     true,
	"static-version":   true,
	"local-file":       true,
}

var engineVendors = map[string]map[string]bool{
	"vllm":         {"nvidia": true, "amd": true, "intel": true},
	"sglang":       {"nvidia": true, "amd": true, "intel": true},
	"llama-cpp":    {"cpu": true, "nvidia": true, "amd": true, "intel": true},
	"tensorrt-llm": {"nvidia": true},
}

var allowedCertifications = map[string]bool{
	"hardware-certified": true,
	"hardware-untested":  true,
	"cpu-only":           true,
}

func ValidateManifest(m ImageManifest) error {
	var errs []error
	if m.Schema != ManifestSchema {
		errs = append(errs, fmt.Errorf("schema must be %q", ManifestSchema))
	}
	if !flavorRe.MatchString(m.Name) {
		errs = append(errs, fmt.Errorf("invalid flavor name %q", m.Name))
	}
	vendors, ok := engineVendors[m.Engine]
	if !ok {
		errs = append(errs, fmt.Errorf("unsupported engine %q", m.Engine))
	} else if !vendors[m.Vendor] {
		errs = append(errs, fmt.Errorf("engine %q does not support vendor %q", m.Engine, m.Vendor))
	}
	if m.Engine != "" && m.Vendor != "" && m.Name != m.Engine+"-"+m.Vendor && !(strings.HasPrefix(m.Name, m.Engine+"-"+m.Vendor+"-")) {
		errs = append(errs, fmt.Errorf("flavor %q must start with engine-vendor prefix %q", m.Name, m.Engine+"-"+m.Vendor))
	}
	if m.Image.Repository != RepositoryPrefix+m.Name {
		errs = append(errs, fmt.Errorf("repository must be %s%s", RepositoryPrefix, m.Name))
	}
	if len(m.Image.Platforms) == 0 {
		errs = append(errs, errors.New("at least one image platform is required"))
	}
	if m.Image.Base.Image == "" || m.Image.Base.Tag == "" {
		errs = append(errs, errors.New("image base image and tag are required"))
	}
	if m.Build.Context == "" || m.Build.Dockerfile == "" {
		errs = append(errs, errors.New("build context and dockerfile are required"))
	}
	if len(m.Tests) == 0 {
		errs = append(errs, errors.New("at least one declared test is required"))
	}
	if !allowedCertifications[m.Publish.Certification] {
		errs = append(errs, fmt.Errorf("unsupported publish certification %q", m.Publish.Certification))
	}
	if (m.Vendor == "amd" || m.Vendor == "intel") && m.Publish.Certification != "hardware-untested" {
		errs = append(errs, fmt.Errorf("%s flavor must be hardware-untested until matching runners exist", m.Vendor))
	}
	if m.Vendor == "cpu" && m.Publish.Certification != "cpu-only" {
		errs = append(errs, errors.New("cpu flavor must use cpu-only publish certification"))
	}
	for _, u := range m.Upstreams {
		if u.Name == "" {
			errs = append(errs, errors.New("upstream name is required"))
		}
		if !allowedUpstreamTypes[u.Type] {
			errs = append(errs, fmt.Errorf("unsupported upstream type %q", u.Type))
		}
	}
	if len(m.Capabilities.Accelerators) == 0 {
		errs = append(errs, errors.New("at least one accelerator is required"))
	} else if m.Vendor == "cpu" && (len(m.Capabilities.Accelerators) != 1 || m.Capabilities.Accelerators[0] != "cpu") {
		errs = append(errs, errors.New("cpu vendor must declare only cpu accelerator"))
	} else if m.Vendor != "cpu" && !contains(m.Capabilities.Accelerators, m.Vendor) {
		errs = append(errs, fmt.Errorf("accelerators must include vendor %q", m.Vendor))
	}
	return errors.Join(errs...)
}

func SupportedUpstreamType(t string) bool { return allowedUpstreamTypes[t] }

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
