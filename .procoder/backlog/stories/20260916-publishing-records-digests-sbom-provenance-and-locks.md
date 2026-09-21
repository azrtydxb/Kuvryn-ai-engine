# Publishing records digests SBOM provenance and locks

Status: done
Created: 2026-09-16
Epic: ci-cd-ct-image-pipeline
Sprint: -

## Description

Publish image flavors by digest only after gates pass, attach flavor-specific tags, and record the exact accepted upstreams, image digest, SBOM, provenance, and test result in the lock/release metadata.

## Acceptance criteria

- [x] Publish workflow attaches immutable engine/vendor/upstream tags.
- [x] `latest` is flavor-specific and never global.
- [x] Published digest is recorded in `upstream.lock`.
- [x] SBOM and provenance artifacts are attached per flavor.
- [x] Failed CT prevents tag movement and lock acceptance.

## Evidence

## Closure evidence

- `go test ./...` passed for schema, detector, CT, placement, runtime wrapper, ops, and publish tag packages.
- `go run ./cmd/kuvryn-manifest-check --root .` validated all 11 engine manifests.
- `procoder test`, `procoder lint`, `procoder security`, and `procoder check` passed with 0 blockers.
- DGX Spark NVIDIA CT evidence was captured under `ct-reports/` for NVIDIA hardware-certified flavors; AMD/Intel remain explicitly build-only and `hardware-untested` because matching hardware is unavailable.
- Accepted upstream locks under `engines/**/upstream.lock` record published digests, CT report paths, and vendor-before-engine publish tags for CPU/NVIDIA certified flavors.
