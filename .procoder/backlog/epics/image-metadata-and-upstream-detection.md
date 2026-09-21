# image-metadata-and-upstream-detection

Status: done
Created: 2026-09-16
Milestone: engine-runtime-foundation-m1
Spec: engine-runtime-foundation

## Description

Define the manifest, lock, and Go detector foundation that lets Kuvryn AI Engine know exactly which engine/vendor image flavors changed and which CI jobs to run. This epic owns schema validation, upstream resolvers, matrix generation, and deterministic affected-flavor planning.

## Closure evidence

- Child stories/tasks have been reconciled against the implemented repository artifacts.
- Final verification commands passed: `go test ./...`, `go run ./cmd/kuvryn-manifest-check --root .`, `procoder test`, `procoder lint`, `procoder security`, and `procoder check`.
- NVIDIA/CPU CT evidence and accepted upstream locks are present for certified/publishable flavors; AMD/Intel limitations are explicitly documented as hardware-untested.
