# Task 6: distributed CT and collectives

Status: done
Created: 2026-09-16
Plan: .procoder/plans/engine-runtime-foundation.md

## Description

Represent distributed CT and collective backend validation for distributed-capable flavors.

## Acceptance criteria

- [x] Two-rank distributed CT fixture exists.
- [x] NVIDIA/NCCL path is wired for distributed-capable flavors on NVIDIA runners.
- [x] AMD/RCCL path is represented but explicitly hardware-untested because no AMD hardware is available.
- [x] Intel/HCCL or oneCCL path is represented but explicitly hardware-untested because no Intel hardware is available.
- [x] Limitations are recorded when distributed support is not yet proven.

## Evidence

- `ci/container-tests/distributed-smoke.sh` requires `WORLD_SIZE>=2`, `RANKS_READY==WORLD_SIZE`, a declared `COLLECTIVE_BACKEND`, and either an OpenAI-compatible inference smoke endpoint or explicit orchestrator-attested inference evidence before it emits `kuvryn.distributed-test-report/v1`.
- `internal/ct/distributed.go` and `internal/ct/validate.go` normalize distributed evidence and enforce required declared tests.
- `internal/placement/collectives.go` represents NCCL, RCCL, HCCL, and oneCCL environment planning.
- `.github/workflows/engine-nvidia-ct.yml` and `ci/container-tests/nvidia-ct-report.sh` provide the NVIDIA-only hardware CT path.
- AMD/Intel distributed CT is intentionally out of local scope: those flavors remain build-only and hardware-untested because no AMD/Intel hardware is available.
- DGX Spark `dgx-spark2` pulled KW-built `llama-cpp-nvidia:kw-test`; `nvidia-image-check.sh` passed static, NVIDIA device, wrapper, and runtime-command checks.
- DGX Spark endpoint smoke passed for `llama-cpp-nvidia:kw-test` with tiny GGUF fixture using `KUVRYN_SMOKE_PORT=8010`.
