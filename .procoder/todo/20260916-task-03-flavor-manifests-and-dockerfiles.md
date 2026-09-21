# Task 3: flavor manifests and Dockerfiles

Status: done
Created: 2026-09-16
Plan: .procoder/plans/engine-runtime-foundation.md

## Description

Create the milestone 1 engine/vendor flavor manifests and lean Dockerfiles that install the runtime command for each image flavor before any build/CT pipeline spends cluster resources.

## Acceptance criteria

- [x] Full target flavor matrix directories exist under `engines/`.
- [x] Each flavor has `image.yaml` and `Dockerfile`.
- [x] Dockerfiles install or verify their engine runtime command (`llama-server`, `vllm`, `sglang`, or `trtllm-serve`) instead of publishing placeholder base images.
- [x] Ray is declared only where the supported engine mode needs it.
- [x] Disallowed all-in-one/all-vendor/`*-ray` flavors do not exist.
- [x] Schema validation passes for every manifest.

## Evidence

- `go run ./cmd/kuvryn-manifest-check --root .` validates 11 manifests.
- Manifests cover vLLM, SGLang, llama.cpp, and TensorRT-LLM across the milestone 1 vendor set.
- `internal/manifestcheck` rejects disallowed aggregate flavors and placeholder pins.
- `internal/manifestcheck.TestEngineDockerfilesInstallRuntimeCommands` rejects Dockerfiles that do not install/verify their runtime command.
- DGX Spark `dgx-spark2` pulled KW-built `llama-cpp-nvidia:kw-test`; `nvidia-image-check.sh` passed static, NVIDIA device, wrapper, and runtime-command checks.
- DGX Spark endpoint smoke passed for `llama-cpp-nvidia:kw-test` with tiny GGUF fixture using `KUVRYN_SMOKE_PORT=8010`.
