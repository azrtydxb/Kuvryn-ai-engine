# Lean Dockerfiles build only required stacks

Status: done
Created: 2026-09-16
Epic: lean-engine-flavor-builds
Sprint: -

## Description

Create lean Dockerfiles for each flavor that install only the engine and matching vendor stack required by that flavor, and verify the runtime command exists during image construction.

## Acceptance criteria

- [x] Every flavor directory has a Dockerfile.
- [x] NVIDIA Dockerfiles do not install ROCm or Intel oneAPI stacks.
- [x] AMD Dockerfiles do not install CUDA or Intel oneAPI stacks.
- [x] Intel Dockerfiles do not install CUDA or ROCm stacks.
- [x] llama.cpp CPU Dockerfile does not install GPU vendor stacks.
- [x] Static CT can inspect the image labels, expected commands, and wrong-stack absence.

## Evidence

## Closure evidence

- `go test ./...` passed for schema, detector, CT, placement, runtime wrapper, ops, and publish tag packages.
- `go run ./cmd/kuvryn-manifest-check --root .` validated all 11 engine manifests.
- `procoder test`, `procoder lint`, `procoder security`, and `procoder check` passed with 0 blockers.
- DGX Spark NVIDIA CT evidence was captured under `ct-reports/` for NVIDIA hardware-certified flavors; AMD/Intel remain explicitly build-only and `hardware-untested` because matching hardware is unavailable.
- Accepted upstream locks under `engines/**/upstream.lock` record published digests, CT report paths, and vendor-before-engine publish tags for CPU/NVIDIA certified flavors.
