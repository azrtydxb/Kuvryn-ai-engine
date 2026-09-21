# Gap analysis and completion plan

This document records the current gap analysis for the engine runtime foundation and the completion lanes that can be run independently.

## Evidence ledger

- Read-only reports: `procoder test`, `procoder lint`, `procoder security`, and `procoder check` report 0 blockers.

## Implementation gaps closed in this pass

- KW Kubernetes nodes now trust the Nexus registry push endpoint `192.168.10.131:5000` through containerd `certs.d` host config.
- `ci/kw/buildctl-build-image.sh` provides a repeatable KW BuildKit build path using the client TLS secret material in `.buildkit-tls/`, cross-compiling `kuvryn-runtime-wrapper` for the target platform before streaming the context.
- `ci/kubernetes/fix-k3s-registry-ca.sh` records the K3s/containerd node CA repair used for Nexus pull validation.
- `ci/kubernetes/pull-smoke.sh` validates image pullability from every KW node with `imagePullPolicy: Always`.
- `ci/kubernetes/engine-runtime-smoke-job.sh` runs a Kubernetes-native OpenAI-compatible smoke job for built engine images without relying on local Docker.

## External/runtime evidence recorded

- KW BuildKit produced immutable CPU/NVIDIA test digests for the publishable local scope.
- CPU image checks passed for `llama-cpp-cpu`.
- NVIDIA hardware CT and endpoint smoke passed on `dgx-spark2` for `llama-cpp-nvidia`, `vllm-nvidia`, `sglang-nvidia`, and `tensorrt-llm-nvidia`.
- TensorRT-LLM endpoint smoke uses a tiny local fp16 Llama HF checkpoint fixture with head size 64 and `KUVRYN_TRTLLM_SERVE_ARGS` targeting `trtllm-serve serve ... --backend pytorch`.
- Accepted `upstream.lock` files record final digests, CT report paths, and vendor-before-engine publish tags for CPU/NVIDIA certified flavors.
- GHCR workflow definitions record SBOM/provenance artifact handling and fail closed on missing required CT evidence.

## Repeatable validation lanes

1. **Build lane**: run `ci/kw/buildctl-build-image.sh <flavor> <ref>` for a selected flavor on the KW BuildKit service.
2. **Kubernetes CT lane**: run `ci/kubernetes/pull-smoke.sh` and `ci/kubernetes/engine-runtime-smoke-job.sh` for Kubernetes-native validation where a fixture exists.
3. **DGX CT lane**: pull KW/GHCR NVIDIA image refs on DGX Spark and run `ci/container-tests/nvidia-image-check.sh` plus `ci/container-tests/engine-runtime-smoke.sh`.
4. **Evidence lane**: validate reports with `cmd/kuvryn-ct-validate` and accept locks with `cmd/kuvryn-lock-accept`.
5. **Publish lane**: after CT evidence is present, run the GHCR publish workflow or `ci/build/publish-image.sh` with immutable digests and vendor-before-engine tags.

## Known policies

- KW cluster builds images.
- DGX Spark hosts validate NVIDIA runtime/hardware CT only; they do not build images.
- AMD and Intel remain build-only and `hardware-untested` until matching hardware exists.
- Do not move `latest` for `hardware-untested` images.
