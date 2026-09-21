# Topology and RDMA validation feed placement

Status: done
Created: 2026-09-16
Epic: cluster-network-rdma-discovery
Sprint: -

## Description

Classify cluster topology and validate RDMA/fabric health so distributed and RDMA-required workloads fail closed when the network cannot safely support them.

## Acceptance criteria

- [x] Classifier handles single-node, direct, ring, switch, Kubernetes fabric, ethernet-only, and unknown.
- [x] RDMA validation distinguishes IP reachability from real RDMA transfer capability when tools are available.
- [x] Kubernetes validation distinguishes node-visible from pod-visible RDMA.
- [x] RDMA-required plans are refused when validation is missing, failed, stale, or pod-invisible.
- [x] Tests cover all topology classes and failure paths.

## Evidence

## Closure evidence

- `go test ./...` passed for schema, detector, CT, placement, runtime wrapper, ops, and publish tag packages.
- `go run ./cmd/kuvryn-manifest-check --root .` validated all 11 engine manifests.
- `procoder test`, `procoder lint`, `procoder security`, and `procoder check` passed with 0 blockers.
- DGX Spark NVIDIA CT evidence was captured under `ct-reports/` for NVIDIA hardware-certified flavors; AMD/Intel remain explicitly build-only and `hardware-untested` because matching hardware is unavailable.
- Accepted upstream locks under `engines/**/upstream.lock` record published digests, CT report paths, and vendor-before-engine publish tags for CPU/NVIDIA certified flavors.
