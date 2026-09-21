# Collective environment planner supports NVIDIA AMD Intel

Status: done
Created: 2026-09-16
Epic: placement-runtime-integration
Sprint: -

## Description

Generate collective environment plans from live discovery for NCCL, RCCL, HCCL/oneCCL, UCX, and no-RDMA cases, and prevent recipes from baking host-specific interface names by default.

## Acceptance criteria

- [x] NCCL env planning uses validated NVIDIA fabric facts.
- [x] RCCL env planning uses validated AMD fabric facts.
- [x] HCCL/oneCCL env planning uses validated Intel fabric facts where supported.
- [x] No-RDMA plans are explicit and refused for RDMA-required deployments.
- [x] Host-specific HCA/interface overrides require an explicit advanced/unsafe setting.
- [x] Tests cover all vendor and no-RDMA paths.

## Evidence

## Closure evidence

- `go test ./...` passed for schema, detector, CT, placement, runtime wrapper, ops, and publish tag packages.
- `go run ./cmd/kuvryn-manifest-check --root .` validated all 11 engine manifests.
- `procoder test`, `procoder lint`, `procoder security`, and `procoder check` passed with 0 blockers.
- DGX Spark NVIDIA CT evidence was captured under `ct-reports/` for NVIDIA hardware-certified flavors; AMD/Intel remain explicitly build-only and `hardware-untested` because matching hardware is unavailable.
- Accepted upstream locks under `engines/**/upstream.lock` record published digests, CT report paths, and vendor-before-engine publish tags for CPU/NVIDIA certified flavors.
