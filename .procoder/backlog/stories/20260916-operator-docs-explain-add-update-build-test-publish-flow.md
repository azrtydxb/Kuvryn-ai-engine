# Operator docs explain add update build test publish flow

Status: done
Created: 2026-09-16
Epic: docs-release-and-operator-experience
Sprint: -

## Description

Document the workflow for adding a new flavor, detecting upstream drift, building, testing, publishing, and diagnosing failures.

## Acceptance criteria

- [x] README explains the flavor model and no-all-in-one rule.
- [x] README or docs show how to add a manifest and Dockerfile.
- [x] Docs show detector command examples and output shape.
- [x] Docs show CI/CT/publish flow and required evidence.
- [x] Docs link image pipeline, flavor matrix, detector, CT, network/RDMA, and operability designs.

## Evidence

## Closure evidence

- `go test ./...` passed for schema, detector, CT, placement, runtime wrapper, ops, and publish tag packages.
- `go run ./cmd/kuvryn-manifest-check --root .` validated all 11 engine manifests.
- `procoder test`, `procoder lint`, `procoder security`, and `procoder check` passed with 0 blockers.
- DGX Spark NVIDIA CT evidence was captured under `ct-reports/` for NVIDIA hardware-certified flavors; AMD/Intel remain explicitly build-only and `hardware-untested` because matching hardware is unavailable.
- Accepted upstream locks under `engines/**/upstream.lock` record published digests, CT report paths, and vendor-before-engine publish tags for CPU/NVIDIA certified flavors.
