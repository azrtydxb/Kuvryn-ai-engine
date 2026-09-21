# Placement selects compatible flavor or fails closed

Status: done
Created: 2026-09-16
Epic: placement-runtime-integration
Sprint: -

## Description

Connect deployment requests, flavor compatibility keys, hardware inventory, topology, and fabric validation to image selection and placement planning.

## Acceptance criteria

- [x] Planner selects image flavor by engine, vendor, accelerator stack, collective, and distributed mode.
- [x] Mixed-vendor auto-placement fails without explicit layout.
- [x] GGUF/llama.cpp, vLLM/SGLang, and TensorRT-LLM compatibility rules are represented in planner tests.
- [x] Planner refuses silent fallback to another vendor or CPU.
- [x] Error messages name the missing capability or incompatible selection.

## Evidence

## Closure evidence

- `go test ./...` passed for schema, detector, CT, placement, runtime wrapper, ops, and publish tag packages.
- `go run ./cmd/kuvryn-manifest-check --root .` validated all 11 engine manifests.
- `procoder test`, `procoder lint`, `procoder security`, and `procoder check` passed with 0 blockers.
- DGX Spark NVIDIA CT evidence was captured under `ct-reports/` for NVIDIA hardware-certified flavors; AMD/Intel remain explicitly build-only and `hardware-untested` because matching hardware is unavailable.
- Accepted upstream locks under `engines/**/upstream.lock` record published digests, CT report paths, and vendor-before-engine publish tags for CPU/NVIDIA certified flavors.
