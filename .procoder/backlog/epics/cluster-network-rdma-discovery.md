# cluster-network-rdma-discovery

Status: done
Created: 2026-09-16
Milestone: engine-runtime-foundation-m1
Spec: engine-runtime-foundation

## Description

Implement Sparkrun-class cluster discovery generalized for NVIDIA, AMD, Intel, Kubernetes, and direct execution. This epic owns host probes, topology classification, RDMA validation, pod-level RDMA checks, collective environment planning inputs, and fail-closed placement facts.

## Closure evidence

- Child stories/tasks have been reconciled against the implemented repository artifacts.
- Final verification commands passed: `go test ./...`, `go run ./cmd/kuvryn-manifest-check --root .`, `procoder test`, `procoder lint`, `procoder security`, and `procoder check`.
- NVIDIA/CPU CT evidence and accepted upstream locks are present for certified/publishable flavors; AMD/Intel limitations are explicitly documented as hardware-untested.
