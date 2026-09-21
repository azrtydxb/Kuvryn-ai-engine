package ct

import (
	"testing"
	"time"
)

func TestReportValidate(t *testing.T) {
	r := Report{Schema: ReportSchema, Flavor: "llama-cpp-cpu", Image: "ghcr.io/azrtydxb/kuvryn-ai-engine/llama-cpp-cpu", StartedAt: time.Now(), EndedAt: time.Now(), Results: []Result{{Name: "static", Type: "static", Status: "passed", Required: true}}}
	if err := r.Validate(); err != nil {
		t.Fatal(err)
	}
	if _, err := r.JSON(); err != nil {
		t.Fatal(err)
	}
}

func TestValidateDistributedResult(t *testing.T) {
	if err := ValidateDistributedResult(DistributedResult{Flavor: "vllm-nvidia", WorldSize: 2, RanksReady: 2, Backend: "nccl", Inference: true}); err != nil {
		t.Fatal(err)
	}
	if err := ValidateDistributedResult(DistributedResult{WorldSize: 1}); err == nil {
		t.Fatal("expected error for one rank")
	}
}
