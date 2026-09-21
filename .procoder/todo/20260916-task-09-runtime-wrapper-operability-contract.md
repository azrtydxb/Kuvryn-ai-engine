# Task 9: runtime wrapper operability contract

Status: done
Created: 2026-09-16
Plan: .procoder/plans/engine-runtime-foundation.md

## Description

Provide a runtime wrapper and unified health/lifecycle contract for engine images.

## Acceptance criteria

- [x] Wrapper starts one runtime process and exits with its status.
- [x] Wrapper forwards SIGTERM/SIGINT and reaps children.
- [x] Health endpoints or adapters exist for live, ready, and startup.
- [x] Lifecycle state and NDJSON events are written.
- [x] Tests cover startup failure, graceful shutdown, child reaping, and readiness semantics.

## Evidence

- `cmd/kuvryn-runtime-wrapper/main.go` starts the configured runtime command, forwards SIGTERM/SIGINT, waits for exit, kills after timeout, writes lifecycle state, and appends lifecycle NDJSON events.
- `internal/runtimecontract/health.go` exposes live/ready/startup handlers.
- `internal/runtimecontract/state.go` writes lifecycle state and NDJSON runtime events; `internal/runtimecontract/metrics.go` writes Prometheus-style metrics.
- `internal/runtimecontract/runtimecontract_test.go` covers health response semantics, state/event writing, and metrics helpers.
- `cmd/kuvryn-runtime-wrapper/main_test.go` builds the wrapper and covers startup failure, child exit status propagation, readiness transition, graceful SIGTERM forwarding, child reap, final state, and lifecycle events.
- `go test ./...` passes.
