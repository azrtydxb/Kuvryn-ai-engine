# Manifest and lock schemas validate engine image flavors

Status: done
Created: 2026-09-16
Epic: image-metadata-and-upstream-detection
Sprint: -

## Description

Implement the `kuvryn.engine-image/v1` and `kuvryn.engine-lock/v1` schema model so each flavor has a strict, testable contract for identity, upstreams, capabilities, build inputs, tests, and publish policy.

## Acceptance criteria

- [x] `internal/engineimage` defines manifest and lock structs.
- [x] Loader rejects unknown fields and malformed YAML.
- [x] Validation enforces flavor naming, engine/vendor compatibility, namespace `ghcr.io/azrtydxb/kuvryn-ai-engine`, and at least one declared test.
- [x] Tests cover valid manifest, missing required fields, wrong namespace, unsupported upstream type, and malformed lock.

## Evidence

## Closure evidence

- `go test ./...` passed for schema, detector, CT, placement, runtime wrapper, ops, and publish tag packages.
- `go run ./cmd/kuvryn-manifest-check --root .` validated all 11 engine manifests.
- `procoder test`, `procoder lint`, `procoder security`, and `procoder check` passed with 0 blockers.
- DGX Spark NVIDIA CT evidence was captured under `ct-reports/` for NVIDIA hardware-certified flavors; AMD/Intel remain explicitly build-only and `hardware-untested` because matching hardware is unavailable.
- Accepted upstream locks under `engines/**/upstream.lock` record published digests, CT report paths, and vendor-before-engine publish tags for CPU/NVIDIA certified flavors.
