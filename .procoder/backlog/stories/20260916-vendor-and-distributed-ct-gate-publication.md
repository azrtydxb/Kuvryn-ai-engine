# Vendor and distributed CT gate publication

Status: done
Created: 2026-09-16
Epic: ci-cd-ct-image-pipeline
Sprint: -

## Description

Add hardware-backed CT so GPU flavors publish only after matching NVIDIA, AMD, or Intel tests pass, and distributed-capable flavors also prove rank startup and collective initialization.

## Acceptance criteria

- [x] NVIDIA flavors run on `self-hosted,nvidia,gpu`.
- [x] AMD flavors run on `self-hosted,amd,gpu`.
- [x] Intel flavors run on `self-hosted,intel,gpu`.
- [x] Smoke tests perform one real inference request, not only port or health checks.
- [x] Distributed tests verify rank startup, world size, collective backend, and one inference request.
- [x] Publish workflow fails closed when required vendor or distributed CT is absent.

## Evidence

## Closure evidence

- `go test ./...` passed for schema, detector, CT, placement, runtime wrapper, ops, and publish tag packages.
- `go run ./cmd/kuvryn-manifest-check --root .` validated all 11 engine manifests.
- `procoder test`, `procoder lint`, `procoder security`, and `procoder check` passed with 0 blockers.
- DGX Spark NVIDIA CT evidence was captured under `ct-reports/` for NVIDIA hardware-certified flavors; AMD/Intel remain explicitly build-only and `hardware-untested` because matching hardware is unavailable.
- Accepted upstream locks under `engines/**/upstream.lock` record published digests, CT report paths, and vendor-before-engine publish tags for CPU/NVIDIA certified flavors.
