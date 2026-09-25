# Decisions

No open decisions.

Resolved in `docs/engine-image-pipeline-design.md`:

- Registry namespace: `ghcr.io/azrtydxb/kuvryn-ai-engine/<flavor>`.
- Upstream detector language: Go first, Python only if Go has practical API or registry integration issues.
- Hardware CT runner target: NVIDIA, AMD, and Intel.
- Milestone 1 scope: full target flavor set, including TensorRT-LLM/NVIDIA.

Resolved by Pascal on 2026-09-21, before the initial commit:

- Engine images run as root: keep root for now, suppressed inline (`nosemgrep`) in all 11 Dockerfiles with a `debt:` marker to revisit once non-root is verified against GPU device access on the NVIDIA, AMD and Intel CT runners.
- `exec.Command` in `cmd/kuvryn-runtime-wrapper/main.go`: suppressed inline; the wrapper execs the image's ENTRYPOINT/CMD argv and the pod's `KUVRYN_READY_COMMAND` by design.

Resolved by Pascal on 2026-09-25:

- vllm-amd stack: ROCm 10.0.0 + vLLM 0.30.0 built from source with gfx1200/gfx1201 targets, not the official rocm723 wheel and not the nightly-rocm100 image.
- llama.cpp: bump every flavor to v0.5.0; llama-cpp-amd moves to ROCm 10.0.0 with GPU targets that include gfx1200;gfx1201.
- Add a Vulkan llama.cpp flavor (`llama-cpp-amd-vulkan`, variant of the AMD flavor) next to the HIP one; hardware-untested; llama.cpp issue #26663 noted as a known RDNA4 risk.
