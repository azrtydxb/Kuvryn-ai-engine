# Task 5: static and vendor container tests

Status: done
Created: 2026-09-16
Plan: .procoder/plans/engine-runtime-foundation.md

## Description

Add static, device, and inference smoke scripts plus publish evidence validation.

## Acceptance criteria

- [x] Static CT checks labels, wrapper, commands, wrong vendor stacks, and SBOM.
- [x] NVIDIA, AMD, Intel, and CPU device checks exist.
- [x] Smoke tests perform one real inference request where supported.
- [x] CT emits `kuvryn.container-test-report/v1` JSON.
- [x] Release publishing fails closed when required CT is absent or failed.

## Evidence

- `ci/container-tests/static-check.sh` validates flavor/vendor/certification labels, wrapper entrypoint, runtime command, emits static+SBOM results, and includes a CPU `vendor-device` result for CPU-only local/CI builds.
- `ci/container-tests/vendor-device-check.sh` handles CPU, NVIDIA, AMD, and Intel checks.
- `ci/container-tests/nvidia-image-check.sh` runs a direct DGX-safe NVIDIA image check combining static labels, in-container GPU visibility, wrapper execution, and required runtime-command presence evidence.
- `ci/container-tests/cpu-image-check.sh` performs the equivalent CPU static/wrapper/runtime-command check for CPU flavors.
- `ci/container-tests/openai-chat-smoke.sh` sends a real OpenAI-compatible chat request.
- `internal/ct/validate.go` rejects missing/failed required CT and missing vendor-device CT for hardware-certified flavors.
- Manual smoke showed `vllm-nvidia` without vendor CT is rejected and `vllm-amd` static-only evidence is accepted as hardware-untested.
- Local `llama-cpp-cpu` build/test smoke passed: `ci/build/build-image.sh llama-cpp-cpu ...`, `static-check.sh`, `kuvryn-ct-validate`, and wrapper `/bin/true` execution.
- Direct DGX smoke on `dgx-spark2` previously exposed that skeleton NVIDIA images could pass static/device/wrapper checks while missing `llama-server`, `vllm`, and `sglang`; runtime-command evidence is now required so those incomplete images cannot be accepted.
- DGX Spark `dgx-spark2` pulled KW-built `llama-cpp-nvidia:kw-test`; `nvidia-image-check.sh` passed static, NVIDIA device, wrapper, and runtime-command checks.
- DGX Spark endpoint smoke passed for `llama-cpp-nvidia:kw-test` with tiny GGUF fixture using `KUVRYN_SMOKE_PORT=8010`.
