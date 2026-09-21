# Go detector emits affected flavor matrix

Status: done
Created: 2026-09-16
Epic: image-metadata-and-upstream-detection
Sprint: -

## Description

Build the Go-first detector that resolves manifest upstreams, compares them to locks, and emits deterministic CI matrix JSON for only affected flavors.

## Acceptance criteria

- [x] `cmd/kuvryn-engine-detect` scans `engines/**/image.yaml`.
- [x] Resolver support exists for GitHub release/tag/commit, container digest, PyPI version, static version, and local file digest.
- [x] Missing lock marks a flavor changed; malformed lock fails detection.
- [x] Matrix output is sorted and deterministic.
- [x] Tests prove an upstream change in one flavor does not rebuild unrelated flavors.

## Evidence

## Closure evidence

- `go test ./...` passed for schema, detector, CT, placement, runtime wrapper, ops, and publish tag packages.
- `go run ./cmd/kuvryn-manifest-check --root .` validated all 11 engine manifests.
- `procoder test`, `procoder lint`, `procoder security`, and `procoder check` passed with 0 blockers.
- DGX Spark NVIDIA CT evidence was captured under `ct-reports/` for NVIDIA hardware-certified flavors; AMD/Intel remain explicitly build-only and `hardware-untested` because matching hardware is unavailable.
- Accepted upstream locks under `engines/**/upstream.lock` record published digests, CT report paths, and vendor-before-engine publish tags for CPU/NVIDIA certified flavors.
