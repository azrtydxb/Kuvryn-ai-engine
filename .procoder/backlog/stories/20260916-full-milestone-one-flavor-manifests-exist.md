# Full milestone one flavor manifests exist

Status: done
Created: 2026-09-16
Epic: lean-engine-flavor-builds
Sprint: -

## Description

Add image manifests for the full milestone 1 matrix: vLLM NVIDIA/AMD/Intel, SGLang NVIDIA/AMD/Intel, llama.cpp CPU/NVIDIA/AMD/Intel, and TensorRT-LLM NVIDIA.

## Acceptance criteria

- [x] Each target flavor has `engines/<engine>/<vendor>/image.yaml`.
- [x] Each manifest declares engine, vendor, accelerator stack, collective, tests, upstreams, and publish policy.
- [x] Ray appears only in engine/vendor manifests whose supported modes require it.
- [x] No `*-ray`, `all-engines`, `all-gpu-vendors`, or mixed-vendor flavor exists.
- [x] Schema validation passes for every manifest.

## Evidence

## Closure evidence

- `go test ./...` passed for schema, detector, CT, placement, runtime wrapper, ops, and publish tag packages.
- `go run ./cmd/kuvryn-manifest-check --root .` validated all 11 engine manifests.
- `procoder test`, `procoder lint`, `procoder security`, and `procoder check` passed with 0 blockers.
- DGX Spark NVIDIA CT evidence was captured under `ct-reports/` for NVIDIA hardware-certified flavors; AMD/Intel remain explicitly build-only and `hardware-untested` because matching hardware is unavailable.
- Accepted upstream locks under `engines/**/upstream.lock` record published digests, CT report paths, and vendor-before-engine publish tags for CPU/NVIDIA certified flavors.
