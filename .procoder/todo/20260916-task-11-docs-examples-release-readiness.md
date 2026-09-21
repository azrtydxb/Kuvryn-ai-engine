# Task 11: docs examples release readiness

Status: done
Created: 2026-09-16
Plan: .procoder/plans/engine-runtime-foundation.md

## Description

Document operator workflows, provide examples, and verify the foundation through procoder gates.

## Acceptance criteria

- [x] README explains the flavor model and no-all-in-one rule.
- [x] Docs show add/update/detect/build/test/publish flow.
- [x] One CPU and one GPU example flavor exist.
- [x] Example locks show resolved inputs, digest, and test evidence shape.
- [x] Final implementation runs `procoder test`, `procoder review`, and `procoder check`.

## Evidence

- `README.md` describes the flavor model, no aggregate images, quickstart, and example locks.
- `docs/engine-image-pipeline-design.md`, `docs/upstream-detector-design.md`, and `docs/container-test-strategy.md` describe add/update/detect/build/test/publish flows.
- CPU example: `engines/llama-cpp/cpu`.
- GPU example: `engines/vllm/nvidia`.
- Lock examples: `examples/locks/llama-cpp-cpu.upstream.lock` and `examples/locks/vllm-nvidia.upstream.lock`.
- Latest verification: `procoder test`, `procoder lint`, `procoder security`, `procoder review`, and `procoder check` ran with 0 blockers.
- DGX Spark `dgx-spark2` pulled KW-built `llama-cpp-nvidia:kw-test`; `nvidia-image-check.sh` passed static, NVIDIA device, wrapper, and runtime-command checks.
- DGX Spark endpoint smoke passed for `llama-cpp-nvidia:kw-test` with tiny GGUF fixture using `KUVRYN_SMOKE_PORT=8010`.
