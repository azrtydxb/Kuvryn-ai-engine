# ci-cd-ct-image-pipeline

Status: done
Created: 2026-09-16
Milestone: engine-runtime-foundation-m1
Spec: engine-runtime-foundation

## Description

Build the CI/CD/CT pipeline that consumes the affected-flavor matrix, builds only selected images, runs static and hardware-backed container tests, publishes only passing digests, and records SBOM/provenance/test evidence per flavor.

## Closure evidence

- Child stories/tasks have been reconciled against the implemented repository artifacts.
- Final verification commands passed: `go test ./...`, `go run ./cmd/kuvryn-manifest-check --root .`, `procoder test`, `procoder lint`, `procoder security`, and `procoder check`.
- NVIDIA/CPU CT evidence and accepted upstream locks are present for certified/publishable flavors; AMD/Intel limitations are explicitly documented as hardware-untested.
