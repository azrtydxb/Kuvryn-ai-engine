# Task 10: OOM heartbeat debug evidence

Status: done
Created: 2026-09-16
Plan: .procoder/plans/engine-runtime-foundation.md

## Description

Normalize heartbeat payloads, failure classification, OOM evidence, and redacted debug data.

## Acceptance criteria

- [x] Heartbeat payload includes deployment, image, flavor, node/pod/container, state, readiness, and rank/world fields.
- [x] Docker/local OOM evidence parsing is implemented.
- [x] Kubernetes OOM/eviction evidence parsing is implemented.
- [x] Debug bundles redact secrets and avoid environment dumps.
- [x] Tests cover crash, OOM, eviction, health restart, and user stop classifications.

## Evidence

- `internal/ops/heartbeat.go` defines the heartbeat payload.
- `internal/ops/oom.go` classifies container exits and Kubernetes OOMKilled/SystemOOM/Evicted/Unhealthy reasons.
- `internal/ops/debug_bundle.go` redacts secret-like environment entries.
- `internal/ops/ops_test.go` covers OOM, eviction, host OOM, health restart, user stop, and redaction behavior.
