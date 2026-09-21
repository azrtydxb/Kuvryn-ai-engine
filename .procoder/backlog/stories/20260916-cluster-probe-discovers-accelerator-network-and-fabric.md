# Cluster probe discovers accelerator network and fabric

Status: done
Created: 2026-09-16
Epic: cluster-network-rdma-discovery
Sprint: -

## Description

Implement the read-only host capability probe and parser for accelerator inventory, management networking, fabric/RDMA interfaces, runtime prerequisites, and explicit limitations.

## Acceptance criteria

- [x] Probe output uses sentinel-delimited accelerator, network, fabric, and runtime sections.
- [x] Parser rejects malformed or missing section boundaries.
- [x] Host record includes management interface/IP, accelerator inventory, fabric interfaces, RDMA facts, collective support, and limitations.
- [x] Tests cover complete, partial, and malformed probe outputs.
- [x] Probe does not print secrets or full environment dumps.

## Evidence

## Closure evidence

- `go test ./...` passed for schema, detector, CT, placement, runtime wrapper, ops, and publish tag packages.
- `go run ./cmd/kuvryn-manifest-check --root .` validated all 11 engine manifests.
- `procoder test`, `procoder lint`, `procoder security`, and `procoder check` passed with 0 blockers.
- DGX Spark NVIDIA CT evidence was captured under `ct-reports/` for NVIDIA hardware-certified flavors; AMD/Intel remain explicitly build-only and `hardware-untested` because matching hardware is unavailable.
- Accepted upstream locks under `engines/**/upstream.lock` record published digests, CT report paths, and vendor-before-engine publish tags for CPU/NVIDIA certified flavors.
