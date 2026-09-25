# Kuvryn AI Engine

Kuvryn AI Engine builds lean LLM inference runtime images per engine and hardware vendor, then rebuilds only the flavors whose declared upstreams changed.

## Quick start

```bash
go test ./...
go run ./cmd/kuvryn-manifest-check --root .
go run ./cmd/kuvryn-engine-detect --root . --output update-plan.json
jq '.flavors[] | select(.changed)' update-plan.json
```

Build one flavor locally after compiling the runtime wrapper:

```bash
go build -o bin/kuvryn-runtime-wrapper ./cmd/kuvryn-runtime-wrapper
ci/build/build-image.sh llama-cpp-cpu engines/llama-cpp/cpu/image.yaml
# CI sets BUILD_PUSH=1 to push a temporary ci-<sha>-<flavor> tag and record image-digests/<flavor>.digest.
mkdir -p ct-reports
ci/container-tests/static-check.sh ghcr.io/azrtydxb/kuvryn-ai-engine/llama-cpp-cpu:ci llama-cpp-cpu \
  > ct-reports/llama-cpp-cpu.json
```

Build on the KW BuildKit service and validate Kubernetes pull/runtime behavior without local Docker:

```bash
ci/kw/buildctl-build-image.sh llama-cpp-cpu \
  192.168.10.131:5000/azrtydxb/kuvryn-ai-engine/llama-cpp-cpu:kw-test
ci/kubernetes/pull-smoke.sh llama-cpp-cpu \
  192.168.10.131:5000/azrtydxb/kuvryn-ai-engine/llama-cpp-cpu:kw-test
ci/kubernetes/engine-runtime-smoke-job.sh llama-cpp-cpu \
  192.168.10.131:5000/azrtydxb/kuvryn-ai-engine/llama-cpp-cpu:kw-test
```

Validate CT evidence and accept a lock after a successful build/test:

```bash
go run ./cmd/kuvryn-ct-validate \
  --manifest engines/llama-cpp/cpu/image.yaml \
  --report ct-reports/llama-cpp-cpu.json

go run ./cmd/kuvryn-lock-accept \
  --manifest engines/llama-cpp/cpu/image.yaml \
  --report ct-reports/llama-cpp-cpu.json \
  --digest sha256:<published-image-digest>

ci/build/write-publish-artifacts.sh \
  llama-cpp-cpu \
  sha256:<published-image-digest> \
  engines/llama-cpp/cpu/image.yaml \
  ct-reports/llama-cpp-cpu.json
```

## Flavor model

A flavor is `<engine>-<vendor>`, for example `vllm-nvidia`, `sglang-amd`, or `llama-cpp-cpu`. The repository deliberately does not define an all-in-one image, an all-vendor GPU image, or a separate `*-ray` flavor axis. Ray is an upstream dependency inside the engine/vendor images that need it.

Images publish under:

```text
ghcr.io/azrtydxb/kuvryn-ai-engine/<flavor>
```

After build and CT pass, publish tags include the immutable digest tag plus version tags:

- stable engine-version tag: `v<kuvryn-version>_<vendor>_<engine>_<engine-version>`; example `v1.0.2_nvidia_vllm_2.3.3`
- unstable nightly tag: `v<kuvryn-version>_<vendor>_<engine>_nightly`; example `v1.0.2_nvidia_vllm_nightly`

`VERSION` changes only when Kuvryn AI Engine changes. The vendor component is the flavor vendor (`nvidia`, `amd`, `intel`, or `cpu`). The engine-version component changes only when that engine upstream version changes.

## Milestone 1 flavors

- `vllm-nvidia`, `vllm-amd`, `vllm-intel`
- `sglang-nvidia`, `sglang-amd`, `sglang-intel`
- `llama-cpp-cpu`, `llama-cpp-nvidia`, `llama-cpp-amd`, `llama-cpp-amd-vulkan`, `llama-cpp-intel`
- `tensorrt-llm-nvidia`

AMD and Intel flavors are built for now, but marked `hardware-untested` because this project does not have AMD/Intel hardware. They must not be treated as hardware-certified artifacts. Hardware-backed CT is NVIDIA-only in this repository.

## Examples

- `examples/locks/llama-cpp-cpu.upstream.lock` shows CPU-only accepted lock evidence, allowed publish tags, and `latest` movement.
- `examples/locks/vllm-nvidia.upstream.lock` shows hardware-certified GPU lock evidence including vendor CT and allowed publish tags.

## Design docs

- [Engine image pipeline](docs/engine-image-pipeline-design.md)
- [Engine flavor matrix](docs/engine-flavor-matrix.md)
- [Upstream detector](docs/upstream-detector-design.md)
- [Container test strategy](docs/container-test-strategy.md)
- [Cluster network and RDMA discovery](docs/cluster-network-rdma-design.md)
- [Runtime operability tooling](docs/runtime-operability-tooling.md)
- [Gap analysis and completion plan](docs/gap-analysis-and-completion-plan.md)
