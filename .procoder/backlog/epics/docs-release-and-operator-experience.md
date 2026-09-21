# docs-release-and-operator-experience

Status: done
Created: 2026-09-16
Milestone: engine-runtime-foundation-m1
Spec: engine-runtime-foundation

## Description

Make the foundation usable by operators and future contributors. This epic owns README flow, example manifests and locks, cross-linked design docs, release evidence, Procoder hygiene, and the documented path for adding or updating a flavor.

## Closure evidence

- Child stories/tasks have been reconciled against the implemented repository artifacts.
- Final verification commands passed: `go test ./...`, `go run ./cmd/kuvryn-manifest-check --root .`, `procoder test`, `procoder lint`, `procoder security`, and `procoder check`.
- NVIDIA/CPU CT evidence and accepted upstream locks are present for certified/publishable flavors; AMD/Intel limitations are explicitly documented as hardware-untested.
