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
