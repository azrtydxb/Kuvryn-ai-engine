package main

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/azrtydxb/kuvryn-ai-engine/internal/runtimecontract"
)

func TestWrapperPropagatesChildExitStatus(t *testing.T) {
	wrapper := buildWrapper(t)
	statePath := filepath.Join(t.TempDir(), "state.json")
	cmd := exec.Command(wrapper, "/bin/sh", "-c", "exit 7")
	cmd.Env = append(os.Environ(), "KUVRYN_STATE_PATH="+statePath, "KUVRYN_HEALTH_ADDR="+freeAddr(t))
	err := cmd.Run()
	if exitCode(err) != 7 {
		t.Fatalf("exit = %v, want 7", err)
	}
	state := readState(t, statePath)
	if state.State != runtimecontract.StateFailed {
		t.Fatalf("state = %#v", state)
	}
}

func TestWrapperWritesFailedStateOnStartupFailure(t *testing.T) {
	wrapper := buildWrapper(t)
	statePath := filepath.Join(t.TempDir(), "state.json")
	cmd := exec.Command(wrapper, filepath.Join(t.TempDir(), "missing-command"))
	cmd.Env = append(os.Environ(), "KUVRYN_STATE_PATH="+statePath, "KUVRYN_HEALTH_ADDR="+freeAddr(t))
	err := cmd.Run()
	if exitCode(err) != 1 {
		t.Fatalf("exit = %v, want 1", err)
	}
	state := readState(t, statePath)
	if state.State != runtimecontract.StateFailed || state.Message == "" {
		t.Fatalf("state = %#v", state)
	}
}

func TestWrapperReadinessAndGracefulShutdown(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix signal semantics")
	}
	wrapper := buildWrapper(t)
	dir := t.TempDir()
	statePath := filepath.Join(dir, "state.json")
	eventsPath := filepath.Join(dir, "events.ndjson")
	readyPath := filepath.Join(dir, "ready")
	termPath := filepath.Join(dir, "term")
	addr := freeAddr(t)
	cmd := exec.Command(wrapper, "/bin/sh", "-c", "trap 'echo term > \"$TERM_FILE\"; exit 0' TERM; while true; do sleep 1; done")
	cmd.Env = append(os.Environ(),
		"KUVRYN_STATE_PATH="+statePath,
		"KUVRYN_EVENTS_PATH="+eventsPath,
		"KUVRYN_HEALTH_ADDR="+addr,
		"KUVRYN_READY_COMMAND=test -f "+readyPath,
		"TERM_FILE="+termPath,
	)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = cmd.Process.Kill(); _, _ = cmd.Process.Wait() }()

	waitHTTP(t, "http://"+addr+"/kuvryn/health/live", http.StatusOK)
	waitHTTP(t, "http://"+addr+"/kuvryn/health/ready", http.StatusServiceUnavailable)
	if err := os.WriteFile(readyPath, []byte("ok"), 0o600); err != nil {
		t.Fatal(err)
	}
	waitHTTP(t, "http://"+addr+"/kuvryn/health/ready", http.StatusOK)
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("wrapper exit = %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("wrapper did not exit after SIGTERM")
	}
	if b, err := os.ReadFile(termPath); err != nil || !strings.Contains(string(b), "term") {
		t.Fatalf("child did not observe SIGTERM: content=%q err=%v", b, err)
	}
	state := readState(t, statePath)
	if state.State != runtimecontract.StateExited {
		t.Fatalf("state = %#v", state)
	}
	events, err := os.ReadFile(eventsPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(events), `"state":"ready"`) || !strings.Contains(string(events), `"state":"stopping"`) {
		t.Fatalf("events = %s", events)
	}
}

func buildWrapper(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "kuvryn-runtime-wrapper")
	cmd := exec.Command("go", "build", "-buildvcs=false", "-o", path, ".")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build wrapper: %v\n%s", err, out)
	}
	return path
}

func freeAddr(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	if err := ln.Close(); err != nil {
		t.Fatal(err)
	}
	return addr
}

func waitHTTP(t *testing.T, url string, status int) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var last error
	for ctx.Err() == nil {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == status {
				return
			}
			last = errors.New(resp.Status)
		} else {
			last = err
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("GET %s did not return %d: last=%v", url, status, last)
}

func readState(t *testing.T, path string) runtimecontract.RuntimeState {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var state runtimecontract.RuntimeState
	if err := json.Unmarshal(b, &state); err != nil {
		t.Fatal(err)
	}
	return state
}
