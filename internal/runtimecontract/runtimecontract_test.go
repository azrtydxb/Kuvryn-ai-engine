package runtimecontract

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestHealthHandler(t *testing.T) {
	h := Handler(Health{Live: func() bool { return true }, Ready: func() bool { return false }, Startup: func() bool { return true }})
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/kuvryn/health/live", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("live code = %d", rr.Code)
	}
	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/kuvryn/health/ready", nil))
	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("ready code = %d", rr.Code)
	}
}

func TestStateEventsAndMetrics(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")
	if err := WriteState(path, RuntimeState{State: StateReady, Ready: true}); err != nil {
		t.Fatal(err)
	}
	events := filepath.Join(dir, "events.ndjson")
	if err := AppendEvent(events, RuntimeEvent{State: StateReady, Ready: true}); err != nil {
		t.Fatal(err)
	}
	eventBytes, err := os.ReadFile(events)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(eventBytes), `"state":"ready"`) {
		t.Fatalf("events = %s", eventBytes)
	}
	var b bytes.Buffer
	if err := WritePrometheus(&b, Metrics{Flavor: "llama-cpp-cpu", StartedAt: time.Now().Add(-time.Second), Ready: true}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(b.String(), "kuvryn_runtime_ready") {
		t.Fatalf("metrics = %s", b.String())
	}
}
