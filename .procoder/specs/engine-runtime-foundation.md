# engine-runtime-foundation

Status: complete

## Problem

Kuvryn AI Engine needs a reproducible foundation for building, testing, publishing, and operating LLM inference runtime images across NVIDIA, AMD, Intel, and CPU targets. Without a spec-first foundation, the project risks drifting into oversized all-in-one images, vendor-specific assumptions, untested distributed modes, and cluster launches that look healthy while RDMA, readiness, or observability are broken.

## Users

- Platform operators need lean images they can patch, roll back, and map to specific engine/vendor deployments.
- Runtime engineers need manifests, locks, and CI/CT gates that rebuild only affected flavors.
- Cluster operators need network/RDMA discovery and fail-closed placement for distributed inference.
- Application teams need consistent health, readiness, metrics, lifecycle state, and OOM evidence regardless of the engine.
- Release owners need provenance, SBOMs, immutable digests, and publish gates per image flavor.

## In scope

- Define per-engine/per-vendor image flavor manifests and lock files.
- Implement Go-first upstream detection and affected matrix generation.
- Build and publish full milestone 1 image flavors under `ghcr.io/azrtydxb/kuvryn-ai-engine/<flavor>`.
- Support milestone 1 engines: vLLM, SGLang, llama.cpp, and TensorRT-LLM.
- Support milestone 1 targets: NVIDIA, AMD, Intel, and CPU where the engine applies.
- Include Ray inside engine/vendor images only when a supported engine mode needs it, especially vLLM distributed modes.
- Add cluster network and RDMA discovery for SSH/direct and Kubernetes environments.
- Add NCCL, RCCL, HCCL/oneCCL, UCX, and no-RDMA planning surfaces as applicable.
- Add static, CPU, vendor, distributed, and benchmark container test gates.
- Add a lightweight runtime operability contract for health, readiness, metrics, heartbeat, OOM evidence, lifecycle events, and debug bundles.
- Keep broader exporters, agents, and monitoring systems outside engine images unless a runtime-local shim is required.

## Out of scope

- No all-in-one image containing all engines and all accelerator stacks.
- No separate `*-ray` flavor axis; Ray is an engine dependency when needed.
- No full systemd in ordinary engine containers.
- No Grafana, Prometheus server, Alertmanager, databases, queues, or SSH server inside engine images.
- No automatic host networking mutation during read-only discovery.
- No silent fallback from RDMA-required distributed plans to slow Ethernet.
- No current ASR/TTS runtime implementation.
- No first-milestone requirement for fleet-scale 100-agent acceptance; that belongs to later fleet-management work.

## Constraints

- Images must include only required runtime/vendor pieces to minimize size and update blast radius.
- Upstream changes must rebuild only the flavors that declare those upstreams.
- Missing, malformed, or unsupported upstream metadata fails closed.
- Hardware-certified release publishing requires passing declared CT for the flavor, including matching vendor hardware when the flavor requires it.
- AMD and Intel flavors are built as hardware-untested until matching GPUs/runners are available; they must not be represented as hardware-certified.
- Host-specific HCA/interface names are generated from discovery and must not be baked into recipes.
- Kubernetes RDMA validation must prove pod-level visibility, not only node-level visibility.
- Internal diagnostics must not print tokens, full environment dumps, model registry credentials, or cluster credentials.
- The first detector implementation is Go unless a practical registry/API blocker requires Python.

## Interfaces

- `engines/<engine>/<vendor>/image.yaml` declares image flavor metadata, upstreams, build inputs, capabilities, tests, and publish policy.
- `engines/<engine>/<vendor>/upstream.lock` records last accepted resolved inputs, image digest, allowed publish tags, and test evidence.
- Go detector command emits `kuvryn.engine-update-plan/v1` JSON matrix for CI.
- Container tests emit `kuvryn.container-test-report/v1` JSON.
- Runtime deployments expose `/kuvryn/health/live`, `/kuvryn/health/ready`, `/kuvryn/health/startup`, and `/metrics` or equivalent adapter endpoints.
- Runtime wrappers write `/var/run/kuvryn/runtime-state.json` and `/var/log/kuvryn/events.ndjson`.
- Cluster discovery emits versioned host capability records covering accelerator, management network, fabric, runtime prerequisites, limitations, and collective support.

## Data

- Image manifests live beside the Dockerfile they describe.
- Lock files live beside the manifest and are updated only after build and required CT pass.
- SBOM/provenance attach to each published image digest.
- Cluster discovery records are runtime state, not recipes.
- Test reports are CI artifacts and may be summarized into release metadata.
- Runtime lifecycle events are NDJSON, scoped to the deployment and safe for collection.

## Edge cases

- A flavor has no lock file yet.
- A base image tag points to a new digest but the engine version is unchanged.
- A shared CT script changes and only declared consuming flavors should rebuild.
- A vLLM image needs Ray while llama.cpp and TensorRT-LLM do not.
- A host has GPUs but no RDMA fabric.
- A Kubernetes node exposes RDMA on the host but the runtime pod cannot see `/dev/infiniband` or equivalent devices.
- A distributed plan spans mixed vendors without an explicit layout.
- An engine health endpoint reports healthy before the model can generate tokens.
- The runtime is killed by cgroup OOM, host OOM, Kubernetes eviction, or a driver crash.
- Intel distributed support differs between Gaudi/HCCL and oneAPI/oneCCL paths.

## Failure modes

- Upstream API unavailable: detection fails unless a workflow explicitly chooses offline lock-only comparison.
- Unsupported upstream type: plan generation fails.
- Changed flavor has no declared tests: plan generation fails.
- Missing hardware runner for hardware-certified release CT: certification waits or fails; it is not a pass. AMD and Intel may still produce hardware-untested build artifacts under the explicit deferred policy.
- RDMA validation missing or failed for an RDMA-required plan: placement fails with diagnostics.
- Wrong vendor stack appears in an image: static CT fails.
- Runtime readiness cannot complete real inference: publish/deploy gate fails.
- OOM evidence is unavailable: the incident records an explicit evidence limitation rather than guessing.

## Acceptance criteria

- [x] `image.yaml` and `upstream.lock` schemas are implemented with validation tests for valid, missing, malformed, and unsupported fields.
- [x] The Go upstream detector emits a deterministic affected-flavor matrix and unit tests prove unrelated flavors are not rebuilt.
- [x] Full milestone 1 image manifests exist for vLLM, SGLang, llama.cpp, and TensorRT-LLM across the target vendor matrix.
- [x] CI builds only detector-selected flavors and publishes only by digest after tests pass.
- [x] Static CT verifies labels, expected commands, runtime wrapper contract, no wrong vendor stack, SBOM generation, and image manifest identity.
- [x] Vendor CT runs on available NVIDIA hardware labels and fails closed for hardware-certified publishing when required NVIDIA validation is absent; AMD and Intel builds are marked hardware-untested because no AMD/Intel hardware is available.
- [x] Distributed CT validates rank startup, collective backend initialization, world size agreement, and one inference request for NVIDIA distributed-capable flavors; AMD/Intel distributed paths are represented only as hardware-untested build metadata.
- [x] Cluster discovery parses complete, partial, and malformed probe outputs and records explicit limitations.
- [x] Topology classification covers single-node, direct, ring, switch, Kubernetes fabric, ethernet-only, and unknown.
- [x] RDMA-required workloads are refused when fabric validation is missing, failed, stale, or not visible in the runtime pod.
- [x] Runtime wrapper exposes the unified health/readiness/metrics/lifecycle contract and handles signals without full systemd.
- [x] OOM/failure evidence distinguishes runtime crash, driver crash, cgroup OOM, host OOM, Kubernetes eviction, health-check restart, and user stop where the platform exposes evidence.
- [x] Ray is declared and tested as an upstream input only in engine/vendor images whose supported modes require it; no `*-ray` flavor is produced.
- [x] Documentation links the image pipeline, flavor matrix, detector, CT, network/RDMA, and operability designs.
- [x] `procoder check` and the repository test suite pass once a recognized test setup exists.

## Open questions

## Decisions

- Registry namespace: `ghcr.io/azrtydxb/kuvryn-ai-engine/<flavor>`.
- Detector language: Go first, Python only for practical registry/API blockers.
- Hardware CT target for this repository: NVIDIA only. AMD and Intel remain build-only/hardware-untested because no AMD/Intel hardware is available.
- Milestone 1 scope: full target flavor set, including TensorRT-LLM/NVIDIA.
- Ray policy: include Ray inside engine/vendor images when a supported mode needs it; do not create `*-ray` flavors.
- Operability policy: use minimal supervisor behavior in entrypoints; do not run full systemd in ordinary engine containers.

## Completion evidence

- Implementation artifacts, CT reports, accepted upstream locks, docs, and verification commands satisfy the milestone acceptance criteria.
- Final verification commands: `go test ./...`, `go run ./cmd/kuvryn-manifest-check --root .`, `procoder test`, `procoder lint`, `procoder security`, and `procoder check`.
