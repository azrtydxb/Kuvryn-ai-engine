package ct

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/azrtydxb/kuvryn-ai-engine/internal/engineimage"
)

type EvidenceReport struct {
	Schema  string   `json:"schema"`
	Flavor  string   `json:"flavor"`
	Image   string   `json:"image"`
	Digest  string   `json:"digest,omitempty"`
	Results []Result `json:"results,omitempty"`
	Tests   []Result `json:"tests,omitempty"`
}

func LoadEvidenceReport(path string) (EvidenceReport, error) {
	var report EvidenceReport
	b, err := os.ReadFile(path)
	if err != nil {
		return report, err
	}
	if err := json.Unmarshal(b, &report); err != nil {
		return report, fmt.Errorf("decode CT report %s: %w", path, err)
	}
	return report, nil
}

func ValidatePublishEvidence(manifest engineimage.ImageManifest, report EvidenceReport) error {
	if report.Schema != ReportSchema {
		return fmt.Errorf("CT report schema must be %q", ReportSchema)
	}
	if report.Flavor != manifest.Name {
		return fmt.Errorf("CT report flavor %q does not match manifest flavor %q", report.Flavor, manifest.Name)
	}
	results := report.AllResults()
	if len(results) == 0 {
		return fmt.Errorf("CT report has no results")
	}
	byName := map[string]Result{}
	byType := map[string]Result{}
	for _, result := range results {
		byName[result.Name] = result
		byType[result.Type] = result
	}
	for _, declared := range manifest.Tests {
		if !declared.Required {
			continue
		}
		result, ok := byName[declared.Name]
		if !ok {
			result, ok = byType[declared.Type]
		}
		if !ok {
			return fmt.Errorf("missing required CT result %q", declared.Name)
		}
		if result.Status != "passed" {
			return fmt.Errorf("required CT result %q status is %q", declared.Name, result.Status)
		}
	}
	if manifest.Publish.Certification == "hardware-certified" {
		if err := requirePassedType(manifest, byType, "vendor-device"); err != nil {
			return err
		}
	}
	return nil
}

func (r EvidenceReport) AllResults() []Result {
	if len(r.Results) > 0 {
		return r.Results
	}
	return r.Tests
}

func requirePassedType(manifest engineimage.ImageManifest, results map[string]Result, typ string) error {
	declared := false
	for _, test := range manifest.Tests {
		if test.Type == typ {
			declared = true
			break
		}
	}
	if !declared {
		return fmt.Errorf("hardware-certified flavor %q must declare %s CT", manifest.Name, typ)
	}
	result, ok := results[typ]
	if !ok {
		return fmt.Errorf("hardware-certified flavor %q missing %s CT result", manifest.Name, typ)
	}
	if result.Status != "passed" {
		return fmt.Errorf("hardware-certified flavor %q has %s CT status %q", manifest.Name, typ, result.Status)
	}
	return nil
}
