package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/azrtydxb/kuvryn-ai-engine/internal/runtimecontract"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: kuvryn-runtime-wrapper <command> [args...]")
		os.Exit(2)
	}
	statePath := getenv("KUVRYN_STATE_PATH", "/var/run/kuvryn/runtime-state.json")
	eventsPath := getenv("KUVRYN_EVENTS_PATH", "/var/run/kuvryn/runtime-events.ndjson")
	_ = os.MkdirAll(filepath.Dir(statePath), 0o755)
	_ = os.MkdirAll(filepath.Dir(eventsPath), 0o755)
	var ready atomic.Bool
	// Running the image's ENTRYPOINT/CMD argv is the wrapper's job; it is set at build time, not by users.
	// nosemgrep: go.lang.security.audit.dangerous-exec-command.dangerous-exec-command
	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		writeStateEvent(statePath, eventsPath, runtimecontract.RuntimeState{State: runtimecontract.StateFailed, Message: err.Error()})
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	writeStateEvent(statePath, eventsPath, runtimecontract.RuntimeState{State: runtimecontract.StateStarting, PID: cmd.Process.Pid})
	server := &http.Server{Addr: getenv("KUVRYN_HEALTH_ADDR", ":8088"), Handler: runtimecontract.Handler(runtimecontract.Health{Live: func() bool { return true }, Ready: ready.Load, Startup: func() bool { return true }})}
	go func() { _ = server.ListenAndServe() }()
	if readyCommand := os.Getenv("KUVRYN_READY_COMMAND"); readyCommand != "" {
		go func() {
			for {
				// KUVRYN_READY_COMMAND comes from the pod spec, which already controls the container.
				// nosemgrep: go.lang.security.audit.dangerous-exec-command.dangerous-exec-command
				probe := exec.Command("/bin/sh", "-c", readyCommand)
				if probe.Run() == nil {
					ready.Store(true)
					writeStateEvent(statePath, eventsPath, runtimecontract.RuntimeState{State: runtimecontract.StateReady, PID: cmd.Process.Pid, Ready: true})
					return
				}
				time.Sleep(time.Second)
			}
		}()
	} else {
		ready.Store(true)
		writeStateEvent(statePath, eventsPath, runtimecontract.RuntimeState{State: runtimecontract.StateReady, PID: cmd.Process.Pid, Ready: true})
	}

	sigc := make(chan os.Signal, 2)
	signal.Notify(sigc, syscall.SIGTERM, syscall.SIGINT)
	waitc := make(chan error, 1)
	go func() { waitc <- cmd.Wait() }()
	var err error
	select {
	case sig := <-sigc:
		ready.Store(false)
		writeStateEvent(statePath, eventsPath, runtimecontract.RuntimeState{State: runtimecontract.StateStopping, PID: cmd.Process.Pid, Message: sig.String()})
		_ = cmd.Process.Signal(sig)
		select {
		case err = <-waitc:
		case <-time.After(30 * time.Second):
			_ = cmd.Process.Kill()
			err = <-waitc
		}
	case err = <-waitc:
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
	code := exitCode(err)
	state := runtimecontract.StateExited
	if code != 0 {
		state = runtimecontract.StateFailed
	}
	writeStateEvent(statePath, eventsPath, runtimecontract.RuntimeState{State: state, Message: fmt.Sprint(err)})
	os.Exit(code)
}

func writeStateEvent(statePath, eventsPath string, state runtimecontract.RuntimeState) {
	_ = runtimecontract.WriteState(statePath, state)
	_ = runtimecontract.AppendEvent(eventsPath, runtimecontract.RuntimeEvent{State: state.State, PID: state.PID, Ready: state.Ready, Message: state.Message})
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		return exit.ExitCode()
	}
	return 1
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
