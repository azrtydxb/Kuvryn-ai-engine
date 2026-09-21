# Task 2: Go upstream detector and affected matrix

Status: done
Created: 2026-09-16
Plan: .procoder/plans/engine-runtime-foundation.md

## Description

Implement the Go detector that resolves upstreams, compares locks, and emits deterministic affected-flavor CI matrix JSON.

## Acceptance criteria

- [x] `cmd/kuvryn-engine-detect` scans engine manifests.
- [x] Required resolver types are implemented or fail closed when unsupported.
- [x] Missing lock marks a flavor changed; malformed lock fails.
- [x] Output schema is `kuvryn.engine-update-plan/v1`.
- [x] Tests prove unrelated flavors are not rebuilt.

## Evidence

- `go test ./...` passes, including `internal/detect` tests.
- `go run ./cmd/kuvryn-engine-detect --root . --output /tmp/update-plan.json` emits 11 flavor entries.
- `internal/detect/resolvers.go` fail-closes unsupported upstream types.
