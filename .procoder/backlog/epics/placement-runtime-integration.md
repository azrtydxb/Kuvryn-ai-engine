# placement-runtime-integration

Status: done
Created: 2026-09-16
Milestone: engine-runtime-foundation-m1
Spec: engine-runtime-foundation

## Description

Connect flavor manifests, cluster discovery, topology classification, and collective validation to deployment planning. This epic ensures Kuvryn selects the right engine/vendor image, refuses unsafe mixed-vendor or RDMA-required placements, and emits launch plans with the correct runtime capabilities.

## Closure evidence

- Child stories/tasks have been reconciled against the implemented repository artifacts.
- Final verification commands passed: `go test ./...`, `go run ./cmd/kuvryn-manifest-check --root .`, `procoder test`, `procoder lint`, `procoder security`, and `procoder check`.
- NVIDIA/CPU CT evidence and accepted upstream locks are present for certified/publishable flavors; AMD/Intel limitations are explicitly documented as hardware-untested.
