# Engine image pipeline design

## Goal

Kuvryn AI Engine builds and publishes lean inference runtime images by engine and hardware vendor. It must never require one "everything" image that carries CUDA, ROCm, oneAPI, TensorRT-LLM, vLLM, SGLang, and llama.cpp together.

Each image flavor is independently versioned, rebuilt, tested, published, and rolled back.

## Non-goals

- No all-in-one runtime image.
- No forced rebuild of unrelated engine/vendor flavors when one upstream changes.
- No implicit multi-vendor container runtime. Multi-vendor scheduling belongs in placement and orchestration, not inside a single image.
- No best-effort hardware certification: an image that cannot pass required hardware CT is not described as hardware-certified.

## Related design docs

- [Engine flavor matrix](engine-flavor-matrix.md) defines the full milestone 1 flavor set.
- [Upstream detector design](upstream-detector-design.md) defines Go-first update detection and matrix generation.
- [Container test strategy](container-test-strategy.md) defines static, smoke, vendor, and distributed CT gates.
- [Cluster network and RDMA discovery](cluster-network-rdma-design.md) defines the cluster-building discovery, topology, and fabric-validation requirements.
- [Runtime operability tooling](runtime-operability-tooling.md) defines the shared heartbeat, health, metrics, OOM, supervisor, and diagnostics contract.

## Image flavor model

A flavor is identified by:

```text
<engine>-<vendor>[-<variant>]
```

Examples:

```text
vllm-nvidia
vllm-amd
sglang-nvidia
sglang-intel
llama-cpp-cpu
llama-cpp-nvidia
llama-cpp-amd
tensorrt-llm-nvidia
```

Each flavor owns a manifest:

```text
engines/<engine>/<vendor>/image.yaml
```

The manifest is the single source of truth for build inputs, runtime capability declarations, test expectations, and publish policy.

## Proposed repository layout

```text
engines/
  vllm/
    nvidia/
      Dockerfile
      image.yaml
      upstream.lock
    amd/
      Dockerfile
      image.yaml
      upstream.lock
  sglang/
    nvidia/
      Dockerfile
      image.yaml
      upstream.lock
    intel/
      Dockerfile
      image.yaml
      upstream.lock
  llama-cpp/
    cpu/
      Dockerfile
      image.yaml
      upstream.lock
    nvidia/
      Dockerfile
      image.yaml
      upstream.lock
    amd/
      Dockerfile
      image.yaml
      upstream.lock
  tensorrt-llm/
    nvidia/
      Dockerfile
      image.yaml
      upstream.lock

ci/
  detect-updates/
  build-matrix/
  container-tests/

docs/
  engine-image-pipeline-design.md
```

## Image manifest schema

Example:

```yaml
schema: kuvryn.engine-image/v1
name: vllm-nvidia
engine: vllm
vendor: nvidia
variant: default

image:
  repository: ghcr.io/azrtydxb/kuvryn-ai-engine/vllm-nvidia
  platforms:
    - linux/amd64
  base:
    image: nvcr.io/nvidia/pytorch
    tag: 25.01-py3
    digest: sha256:example

upstreams:
  - name: vllm
    type: github-release
    repo: vllm-project/vllm
    version: v0.10.1
  - name: cuda
    type: container-digest
    image: nvcr.io/nvidia/pytorch
    tag: 25.01-py3
    digest: sha256:example
  - name: python
    type: semver
    version: "3.12"

capabilities:
  accelerator_vendor: nvidia
  accelerator_stack: cuda
  collective: nccl
  distributed: true
  openai_api: true
  readiness:
    health_path: /health
    inference_probe: openai-chat-stream-v1

build:
  dockerfile: Dockerfile
  target: runtime
  args:
    VLLM_VERSION: v0.10.1
  cache_scope: vllm-nvidia

publish:
  certification: hardware-certified
  tags:
    - "{{ engine }}-{{ vendor }}-{{ upstreams.vllm.version }}"
    - "{{ engine }}-{{ vendor }}-latest"

tests:
  smoke:
    command: ci/container-tests/openai-chat-smoke.sh
    requires_gpu: true
  cpu_static:
    command: ci/container-tests/image-static-check.sh
    requires_gpu: false
```

## Lock file model

Each flavor has an `upstream.lock` recording the exact resolved inputs from the last accepted build:

```yaml
schema: kuvryn.engine-lock/v1
flavor: vllm-nvidia
resolvedInputs:
  - name: engine
    type: static-version
    value: 0.29.0
    source: engines/vllm/nvidia/image.yaml
  - name: base-image
    type: container-digest
    value: nvcr.io/nvidia/pytorch:25.01-py3
    digest: sha256:96990c82825613c3bdeebb66675c7c91b0123f64a5895623316dc5b824e0d7a9
    source: engines/vllm/nvidia/image.yaml
imageDigest: sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
publishedTags:
  - bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
  - v0.1.0_nvidia_vllm_0.29.0
  - v0.1.0_nvidia_vllm_nightly
  - latest
  - hardware-certified
testEvidence:
  - name: static
    status: passed
  - name: vendor-device
    status: passed
```

The lock file is updated only after the image builds, `kuvryn-ct-validate` accepts CT evidence, and the published image digest is known. `kuvryn-lock-accept` writes the accepted lock file or proposed lock artifact with the exact publish tags that were allowed for the flavor certification.

## Update detection

`kuvryn-manifest-check` validates all `image.yaml` files before detection. It rejects malformed manifests, missing Dockerfiles, placeholder upstream pins, disallowed aggregate flavors, and accidental hardware-certified AMD/Intel manifests while the deferred CT policy is active.

The detector reads all `image.yaml` files and compares their resolved upstream inputs against `upstream.lock`.

Detection outputs an affected build matrix:

```json
{
  "include": [
    {
      "name": "vllm-nvidia",
      "engine": "vllm",
      "vendor": "nvidia",
      "path": "engines/vllm/nvidia"
    }
  ]
}
```

A flavor is affected when any declared input changed:

- engine release or commit changed
- base image digest changed
- accelerator stack version changed: CUDA, ROCm, oneAPI, Gaudi software
- build script or Dockerfile changed
- shared CI test or shared base layer used by that flavor changed
- manifest changed

A flavor is not affected just because another engine/vendor changed.

## CI/CD/CT flow

```text
Detect upstreams
  -> produce affected flavor matrix
  -> build only affected flavors
  -> run static container checks
  -> run CPU-only smoke where possible
  -> run GPU/vendor CT on matching runners
  -> publish only passing flavors
  -> update lock files / release metadata
```

### 1. Detect

Runs on schedule, manual dispatch, and pull requests that touch engine/build files.

Responsibilities:

- resolve current upstream versions and digests
- compare with `upstream.lock`
- classify flavors as changed or unchanged
- emit GitHub Actions matrix JSON
- never publish

### 2. Build

Builds each affected flavor in isolation.

Rules:

- one matrix entry per image flavor
- scoped cache per flavor, not a shared global cache
- no QEMU for accelerator-heavy builds unless explicitly marked safe
- native runners for architectures that compile vendor libraries
- build by temporary CI tag and record the resulting immutable digest before attaching stable tags
- generate static CT report artifacts and image digest artifacts for every built flavor

### 3. Container tests

Every flavor declares its own test set.

Minimum non-GPU checks:

- image starts
- expected binary imports/loads
- no unrelated accelerator stack present, for example ROCm tools inside `vllm-nvidia`
- expected runtime command exists and is executable; `runtime-command` evidence is required, so a base image wrapper without the engine app installed is not acceptable
- OpenAPI-compatible route can be configured where supported

GPU CT checks run only on available matching hardware runners:

- NVIDIA images on NVIDIA runners
- CPU images on ordinary Linux runners
- AMD and Intel images are build-only/hardware-untested because this project does not have AMD/Intel hardware

Current local policy: NVIDIA is the only hardware-backed CT target. AMD and Intel image flavors are built but hardware-untested because matching GPUs are not available. Their manifests use `publish.certification: hardware-untested` and keep vendor CT declared with `required: false`; release metadata must call them hardware-untested and must not imply AMD/Intel validation.

Minimum GPU CT:

- launch with tiny model or synthetic fixture
- `/health` passes
- one real inference request succeeds
- readiness records time-to-first-token for streaming-capable engines
- distributed flavor validates collective backend wiring when declared

### 4. Publish

Publishing happens only after declared tests pass.

Publish behavior:

- attach immutable digest tags
- attach semantic vendor/engine tags as `v<kuvryn-version>_<vendor>_<engine>_<engine-version>`; for example `v1.0.2_nvidia_vllm_2.3.3`
- attach unstable nightly tags as `v<kuvryn-version>_<vendor>_<engine>_nightly`; for example `v1.0.2_nvidia_vllm_nightly`
- change `VERSION` only when Kuvryn AI Engine changes; change the engine-version tag component only when that engine upstream changes
- attach `latest` only for the same flavor, never globally, and never for `hardware-untested` artifacts
- validate all manifests again before tag movement
- emit SBOM, provenance, and `publish-tags.json` per flavor under `publish-artifacts/<flavor>/`
- retain image digest in release metadata and `upstream.lock`

## Ray policy

Ray is treated as an engine dependency when that engine needs Ray for a supported runtime mode. It does not create a separate image flavor by itself.

Examples:

- `vllm-nvidia` may include Ray when the supported vLLM distributed modes require it.
- `vllm-amd` may include Ray when the supported ROCm vLLM distributed modes require it.
- `llama-cpp-*` and `tensorrt-llm-nvidia` must not include Ray unless their declared engine modes need it.

This keeps the flavor matrix focused on engine and hardware vendor, while still allowing engines to ship the orchestration libraries they actually need.

Ray must still be declared in `image.yaml` as an upstream input, tested in container tests, and included in update detection so Ray changes only rebuild affected engine/vendor images.

## Lean image rules

Each Dockerfile must include only what the flavor requires.

Examples:

- `vllm-nvidia` includes CUDA/NCCL/vLLM, not ROCm, oneAPI, llama.cpp, or SGLang.
- `sglang-intel` includes Intel runtime dependencies, not CUDA or ROCm.
- `llama-cpp-cpu` includes CPU inference dependencies only.
- `tensorrt-llm-nvidia` is NVIDIA-only and separate from `vllm-nvidia`.

Shared base images are allowed only when they are narrow and vendor-specific:

```text
base-nvidia-cuda-runtime
base-amd-rocm-runtime
base-intel-runtime
base-cpu-runtime
```

There is no shared `base-all-accelerators` image.

## Runtime and orchestration contract

The control plane selects images by matching:

- model requirements
- engine family
- GPU vendor
- required capabilities
- distributed/sharding mode
- available cluster hardware

The image declares what it can do; placement decides where it may run.

Example runtime compatibility key:

```yaml
capabilities:
  accelerator_vendor: amd
  accelerator_stack: rocm
  collective: rccl
  distributed: true
```

A deployment requesting distributed vLLM on AMD must resolve to `vllm-amd` and RCCL-capable nodes. It must not fall back to NVIDIA or CPU silently.

## Decisions

1. Container registry namespace: publish under the `azrtydxb/kuvryn-ai-engine` GitHub scope, for example `ghcr.io/azrtydxb/kuvryn-ai-engine/<flavor>`.
2. Upstream detector implementation language: Go first. Python is allowed only if Go creates practical issues for an upstream API or registry integration.
3. Initial vendor CT availability: design for NVIDIA, AMD, and Intel hardware-backed container tests, but build AMD and Intel flavors as hardware-untested until matching GPUs/runners are available.
4. Milestone 1 flavor set: full target set, including TensorRT-LLM/NVIDIA.
