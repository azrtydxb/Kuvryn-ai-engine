# CI builds only detected flavors

Status: done
Created: 2026-09-16
Epic: ci-cd-ct-image-pipeline
Sprint: -

## Description

Wire GitHub Actions so pull requests and scheduled/manual workflows consume the detector matrix and build only affected flavor images.

## Acceptance criteria

- [x] Detection workflow publishes the update-plan artifact.
- [x] Build workflow consumes the matrix include list from the detector.
- [x] Build cache scope is flavor-specific.
- [x] Pull request workflows never publish stable tags.
- [x] CI logs identify skipped unchanged flavors without building them.

## Evidence

## Closure evidence

- `go test ./...` passed for schema, detector, CT, placement, runtime wrapper, ops, and publish tag packages.
- `go run ./cmd/kuvryn-manifest-check --root .` validated all 11 engine manifests.
- `procoder test`, `procoder lint`, `procoder security`, and `procoder check` passed with 0 blockers.
- DGX Spark NVIDIA CT evidence was captured under `ct-reports/` for NVIDIA hardware-certified flavors; AMD/Intel remain explicitly build-only and `hardware-untested` because matching hardware is unavailable.
- Accepted upstream locks under `engines/**/upstream.lock` record published digests, CT report paths, and vendor-before-engine publish tags for CPU/NVIDIA certified flavors.
