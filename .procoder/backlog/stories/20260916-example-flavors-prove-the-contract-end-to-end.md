# Example flavors prove the contract end to end

Status: done
Created: 2026-09-16
Epic: docs-release-and-operator-experience
Sprint: -

## Description

Provide one CPU example and one GPU example that exercise manifest loading, detection, lean Dockerfile runtime installation, static CT, runtime wrapper contract, and documentation flow.

## Acceptance criteria

- [x] Example CPU flavor exists and passes schema/static checks.
- [x] Example GPU flavor exists and passes schema/static checks.
- [x] Example locks demonstrate accepted resolved inputs and digest/test evidence shape.
- [x] Examples are referenced from the README and detector docs.
- [x] Examples do not publish real stable tags by default.

## Evidence

## Closure evidence

- `go test ./...` passed for schema, detector, CT, placement, runtime wrapper, ops, and publish tag packages.
- `go run ./cmd/kuvryn-manifest-check --root .` validated all 11 engine manifests.
- `procoder test`, `procoder lint`, `procoder security`, and `procoder check` passed with 0 blockers.
- DGX Spark NVIDIA CT evidence was captured under `ct-reports/` for NVIDIA hardware-certified flavors; AMD/Intel remain explicitly build-only and `hardware-untested` because matching hardware is unavailable.
- Accepted upstream locks under `engines/**/upstream.lock` record published digests, CT report paths, and vendor-before-engine publish tags for CPU/NVIDIA certified flavors.
