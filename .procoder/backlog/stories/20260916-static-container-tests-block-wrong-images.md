# Static container tests block wrong images

Status: done
Created: 2026-09-16
Epic: ci-cd-ct-image-pipeline
Sprint: -

## Description

Implement static CT that fails before GPU runner time is used when an image has wrong labels, missing commands, missing wrapper contract, wrong vendor stack, or missing SBOM.

## Acceptance criteria

- [x] Static CT verifies required OCI labels and flavor identity.
- [x] Static CT verifies runtime wrapper and expected engine command exist.
- [x] Static CT fails when unrelated vendor stacks are present.
- [x] Static CT verifies SBOM generation succeeds.
- [x] Static CT emits `kuvryn.container-test-report/v1` JSON.

## Evidence

## Closure evidence

- `go test ./...` passed for schema, detector, CT, placement, runtime wrapper, ops, and publish tag packages.
- `go run ./cmd/kuvryn-manifest-check --root .` validated all 11 engine manifests.
- `procoder test`, `procoder lint`, `procoder security`, and `procoder check` passed with 0 blockers.
- DGX Spark NVIDIA CT evidence was captured under `ct-reports/` for NVIDIA hardware-certified flavors; AMD/Intel remain explicitly build-only and `hardware-untested` because matching hardware is unavailable.
- Accepted upstream locks under `engines/**/upstream.lock` record published digests, CT report paths, and vendor-before-engine publish tags for CPU/NVIDIA certified flavors.
