package ct

import (
	"testing"

	"github.com/azrtydxb/kuvryn-ai-engine/internal/engineimage"
)

func TestValidatePublishEvidenceHardwareCertifiedRequiresVendorCT(t *testing.T) {
	manifest := testManifest("vllm-nvidia", "nvidia", "hardware-certified", []engineimage.TestDeclaration{
		{Name: "static", Type: "static", Required: true},
		{Name: "vendor-device", Type: "vendor-device", Required: true},
	})
	report := EvidenceReport{Schema: ReportSchema, Flavor: "vllm-nvidia", Image: "image", Results: []Result{{Name: "static", Type: "static", Status: "passed", Required: true}}}
	if err := ValidatePublishEvidence(manifest, report); err == nil {
		t.Fatal("expected missing vendor CT error")
	}
	report.Results = append(report.Results, Result{Name: "vendor-device", Type: "vendor-device", Status: "passed", Required: true})
	if err := ValidatePublishEvidence(manifest, report); err != nil {
		t.Fatal(err)
	}
}

func TestValidatePublishEvidenceAllowsUntestedOptionalVendorCT(t *testing.T) {
	manifest := testManifest("vllm-amd", "amd", "hardware-untested", []engineimage.TestDeclaration{
		{Name: "static", Type: "static", Required: true},
		{Name: "vendor-device", Type: "vendor-device", Required: false},
	})
	report := EvidenceReport{Schema: ReportSchema, Flavor: "vllm-amd", Image: "image", Results: []Result{{Name: "static", Type: "static", Status: "passed", Required: true}}}
	if err := ValidatePublishEvidence(manifest, report); err != nil {
		t.Fatal(err)
	}
}

func TestValidatePublishEvidenceRejectsFlavorMismatchAndFailedRequired(t *testing.T) {
	manifest := testManifest("llama-cpp-cpu", "cpu", "cpu-only", []engineimage.TestDeclaration{{Name: "static", Type: "static", Required: true}})
	if err := ValidatePublishEvidence(manifest, EvidenceReport{Schema: ReportSchema, Flavor: "other", Image: "image", Results: []Result{{Name: "static", Type: "static", Status: "passed", Required: true}}}); err == nil {
		t.Fatal("expected flavor mismatch error")
	}
	if err := ValidatePublishEvidence(manifest, EvidenceReport{Schema: ReportSchema, Flavor: "llama-cpp-cpu", Image: "image", Results: []Result{{Name: "static", Type: "static", Status: "failed", Required: true}}}); err == nil {
		t.Fatal("expected failed required test error")
	}
}

func testManifest(name, vendor, certification string, tests []engineimage.TestDeclaration) engineimage.ImageManifest {
	return engineimage.ImageManifest{
		Name:    name,
		Engine:  "vllm",
		Vendor:  vendor,
		Publish: engineimage.Publish{Certification: certification},
		Tests:   tests,
	}
}
