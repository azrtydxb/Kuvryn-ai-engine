package ops

import "testing"

func TestFailureClassification(t *testing.T) {
	if got := ClassifyContainerExit(137, false, false).Kind; got != FailureCgroupOOM {
		t.Fatalf("kind = %s", got)
	}
	if got := ClassifyKubernetes("Evicted").Kind; got != FailureEvicted {
		t.Fatalf("kind = %s", got)
	}
	if got := ClassifyContainerExit(1, false, true).Kind; got != FailureUserStop {
		t.Fatalf("kind = %s", got)
	}
	if got := ClassifyKubernetes("SystemOOM").Kind; got != FailureHostOOM {
		t.Fatalf("kind = %s", got)
	}
	if got := ClassifyKubernetes("Unhealthy").Kind; got != FailureHealthRestart {
		t.Fatalf("kind = %s", got)
	}
}

func TestRedactEnv(t *testing.T) {
	got := RedactEnv([]string{"HF_TOKEN=abc", "MODEL=llama"})
	if got[0] != "HF_TOKEN=<redacted>" || got[1] != "MODEL=llama" {
		t.Fatalf("redacted = %#v", got)
	}
}
