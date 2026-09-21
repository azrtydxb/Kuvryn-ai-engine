# runtime-operability-contract

Status: done
Created: 2026-09-16
Milestone: engine-runtime-foundation-m1
Spec: engine-runtime-foundation

## Description

Provide a small common operational surface for every runtime without placing full monitoring stacks or systemd inside engine images. This epic owns the runtime wrapper, health/readiness contract, metrics baseline, heartbeat shape, OOM/failure evidence, lifecycle events, and debug bundle rules.

## Closure evidence

- Child stories/tasks have been reconciled against the implemented repository artifacts.
- Final verification commands passed: `go test ./...`, `go run ./cmd/kuvryn-manifest-check --root .`, `procoder test`, `procoder lint`, `procoder security`, and `procoder check`.
- NVIDIA/CPU CT evidence and accepted upstream locks are present for certified/publishable flavors; AMD/Intel limitations are explicitly documented as hardware-untested.
