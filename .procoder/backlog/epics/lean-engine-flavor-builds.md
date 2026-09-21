# lean-engine-flavor-builds

Status: done
Created: 2026-09-16
Milestone: engine-runtime-foundation-m1
Spec: engine-runtime-foundation

## Description

Create the full milestone 1 engine/vendor image matrix with lean Dockerfiles and image manifests. This epic ensures Ray is included only as a required engine dependency, forbids separate `*-ray` flavors, installs the runtime app before builds are scheduled, and prevents all-in-one or all-vendor images from appearing.

## Closure evidence

- Child stories/tasks have been reconciled against the implemented repository artifacts.
- Final verification commands passed: `go test ./...`, `go run ./cmd/kuvryn-manifest-check --root .`, `procoder test`, `procoder lint`, `procoder security`, and `procoder check`.
- NVIDIA/CPU CT evidence and accepted upstream locks are present for certified/publishable flavors; AMD/Intel limitations are explicitly documented as hardware-untested.
