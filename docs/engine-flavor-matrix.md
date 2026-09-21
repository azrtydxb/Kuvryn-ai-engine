# Engine flavor matrix

## Purpose

This matrix defines the full target set for milestone 1. Each row is an independently built and tested image flavor. A failure or upstream update in one row must not force rebuilding unrelated rows.

Current bootstrap pins are real upstream identifiers, not placeholders: vLLM `0.29.0`, SGLang `0.5.19`, llama.cpp `v0.4.1`, TensorRT-LLM `1.2.1`, and Ray `2.58.0` for Ray-capable vLLM flavors. Base images are recorded with resolved `sha256:` digests in each manifest.

## Naming

Images publish under:

```text
ghcr.io/azrtydxb/kuvryn-ai-engine/<engine>-<vendor>[-<variant>]
```

Examples:

```text
ghcr.io/azrtydxb/kuvryn-ai-engine/vllm-nvidia
ghcr.io/azrtydxb/kuvryn-ai-engine/sglang-intel
ghcr.io/azrtydxb/kuvryn-ai-engine/llama-cpp-cpu
```

## Milestone 1 target flavors

| Flavor                | Engine       | Vendor | Accelerator stack            | Collective                               | Distribution role                                      | Notes                                                                                           |
| --------------------- | ------------ | ------ | ---------------------------- | ---------------------------------------- | ------------------------------------------------------ | ----------------------------------------------------------------------------------------------- |
| `vllm-nvidia`         | vLLM         | NVIDIA | CUDA                         | NCCL                                     | Single-node and distributed                            | Primary CUDA/OpenAI-compatible serving image; include Ray when supported vLLM modes require it. |
| `vllm-amd`            | vLLM         | AMD    | ROCm                         | RCCL                                     | Single-node and distributed                            | ROCm build path; include Ray when supported vLLM modes require it; no CUDA packages.            |
| `vllm-intel`          | vLLM         | Intel  | Intel/Gaudi or oneAPI target | HCCL or oneCCL                           | Single-node first; distributed when upstream is proven | Keep Intel-specific dependencies isolated; include Ray only if required by supported mode.      |
| `sglang-nvidia`       | SGLang       | NVIDIA | CUDA                         | NCCL                                     | Single-node and distributed                            | CUDA SGLang serving image.                                                                      |
| `sglang-amd`          | SGLang       | AMD    | ROCm                         | RCCL                                     | Single-node and distributed when upstream supports it  | ROCm-only; no CUDA packages.                                                                    |
| `sglang-intel`        | SGLang       | Intel  | Intel/Gaudi or oneAPI target | HCCL or oneCCL                           | Single-node first; distributed when upstream is proven | Separate from NVIDIA and AMD images.                                                            |
| `llama-cpp-cpu`       | llama.cpp    | CPU    | CPU                          | none                                     | Single-node                                            | Baseline portable fallback and CI smoke anchor.                                                 |
| `llama-cpp-nvidia`    | llama.cpp    | NVIDIA | CUDA                         | none or NCCL only if supported by mode   | Single-node first                                      | CUDA llama.cpp without vLLM/SGLang.                                                             |
| `llama-cpp-amd`       | llama.cpp    | AMD    | HIP/ROCm                     | none or RCCL only if supported by mode   | Single-node first                                      | HIP llama.cpp without CUDA.                                                                     |
| `llama-cpp-intel`     | llama.cpp    | Intel  | SYCL/oneAPI where supported  | none or oneCCL only if supported by mode | Single-node first                                      | Intel-specific llama.cpp build.                                                                 |
| `tensorrt-llm-nvidia` | TensorRT-LLM | NVIDIA | CUDA + TensorRT-LLM          | NCCL                                     | Single-node and distributed                            | NVIDIA-only. Never bundled into vLLM image.                                                     |

## Optional later variants

Variants are allowed only when they prevent bloat or incompatible dependency sets in the default image. Ray alone is not a variant reason; include it in the engine image when that engine's supported modes require it.

Examples:

| Variant   | Example image           | Reason                                                        |
| --------- | ----------------------- | ------------------------------------------------------------- |
| `nightly` | `vllm-nvidia-nightly`   | Tracks upstream nightly without destabilizing default.        |
| `triton`  | `vllm-nvidia-triton`    | Adds a heavy optional serving/optimization dependency.        |
| `minimal` | `llama-cpp-cpu-minimal` | Strips optional APIs and debug tools.                         |
| `bench`   | `sglang-nvidia-bench`   | Adds benchmark tools that do not belong in production images. |

## Explicit exclusions

These images are intentionally not allowed:

| Disallowed image  | Reason                                                                      |
| ----------------- | --------------------------------------------------------------------------- |
| `all-engines`     | Couples every engine release and creates a large blast radius.              |
| `all-gpu-vendors` | Pulls CUDA, ROCm, and Intel stacks into one image.                          |
| `vllm-sglang-*`   | Runtime families update independently and should be testable independently. |
| `*-ray`           | Ray is an engine dependency when needed, not a separate flavor axis.        |
| `nvidia-amd-*`    | Multi-vendor placement belongs to orchestration, not one container.         |

## Compatibility keys

Each flavor must declare these keys in `image.yaml`:

```yaml
compatibility:
  engine: vllm
  vendor: nvidia
  accelerator_stack: cuda
  collective: nccl
  APIs:
    - openai-chat-completions
  distributed_modes:
    - tensor-parallel
    - pipeline-parallel
```

The scheduler may select a flavor only when the deployment request and node inventory match the declared compatibility keys.
