# engine-runtime-foundation — implementation plan

Status: complete
Spec: .procoder/specs/engine-runtime-foundation.md

## Goal

Build the first shippable Kuvryn AI Engine foundation for lean per-engine/per-vendor inference images, upstream-aware CI/CD/CT, cluster/RDMA discovery, and unified runtime operability.

## Architecture

The foundation is split into narrow packages and pipelines: schemas/manifests define flavor identity, the Go detector computes affected builds, CI builds and tests only affected flavors, cluster discovery produces placement facts, and runtime wrappers expose a uniform operational contract. Engine images stay lean; shared observability and host/fabric agents run as sidecars, DaemonSets, host agents, or control-plane services.

## Constraints

- No all-in-one image and no all-vendor image.
- Publish under `ghcr.io/azrtydxb/kuvryn-ai-engine/<flavor>`.
- Go-first detector; Python fallback only for practical API/registry issues.
- NVIDIA hardware CT is a release gate when NVIDIA hardware is available; AMD and Intel flavors are built as hardware-untested until matching GPUs/runners are available.
- Milestone 1 covers the full target flavor set including TensorRT-LLM/NVIDIA.
- Ray is included only inside engine/vendor images whose supported modes need it; do not create `*-ray` flavors.
- No full systemd in normal engine containers.
- Discovery is read-only unless a user explicitly invokes a setup operation.

## Task 1: repository foundation and schemas

Files:

- `go.mod` / `go.sum`: Go module bootstrap.
- `internal/engineimage/schema.go`: manifest and lock structs.
- `internal/engineimage/validate.go`: validation rules.
- `internal/engineimage/schema_test.go`: valid/invalid schema tests.
- `docs/engine-image-pipeline-design.md`: schema references.

Interfaces:

- `type ImageManifest struct` matching `kuvryn.engine-image/v1`.
- `type UpstreamLock struct` matching `kuvryn.engine-lock/v1`.
- `func LoadManifest(path string) (ImageManifest, error)`.
- `func LoadLock(path string) (UpstreamLock, error)`.
- `func ValidateManifest(m ImageManifest) error`.

Steps:

- [x] Initialize a Go module for `github.com/azrtydxb/kuvryn-ai-engine`.
- [x] Define manifest, upstream, capability, build, publish, test, and lock structs.
- [x] Implement YAML loading with unknown-field rejection.
- [x] Validate flavor naming, engine/vendor compatibility, test presence, and publish target namespace.
- [x] Test missing lock, malformed manifest, unsupported upstream type, invalid flavor, missing test, and valid examples.

## Task 2: Go upstream detector and affected matrix

Files:

- `cmd/kuvryn-engine-detect/main.go`: CLI entrypoint.
- `internal/detect/detect.go`: detector orchestration.
- `internal/detect/resolvers.go`: resolver interfaces.
- `internal/detect/matrix.go`: update-plan output.
- `internal/detect/*_test.go`: resolver and matrix tests.
- `docs/upstream-detector-design.md`: command examples.

Interfaces:

- `kuvryn-engine-detect --root . --output update-plan.json`.
- `type Resolver interface { Resolve(context.Context, Upstream) (ResolvedInput, error) }`.
- Output schema `kuvryn.engine-update-plan/v1`.

Steps:

- [x] Discover `engines/**/image.yaml` files.
- [x] Resolve `github-release`, `github-tag`, `github-commit`, `container-digest`, `pypi-version`, `static-version`, and `local-file` inputs.
- [x] Compare resolved inputs with `upstream.lock`.
- [x] Mark missing lock as changed and malformed lock as failure.
- [x] Emit deterministic JSON sorted by flavor name.
- [x] Test that unrelated flavor changes do not enter the matrix.

## Task 3: milestone 1 flavor manifests and lean Dockerfiles

Files:

- `engines/vllm/{nvidia,amd,intel}/image.yaml` and `Dockerfile`.
- `engines/sglang/{nvidia,amd,intel}/image.yaml` and `Dockerfile`.
- `engines/llama-cpp/{cpu,nvidia,amd,intel}/image.yaml` and `Dockerfile`.
- `engines/tensorrt-llm/nvidia/image.yaml` and `Dockerfile`.
- `docs/engine-flavor-matrix.md`: matrix updates.

Interfaces:

- Each flavor has `image.yaml`, `Dockerfile`, and generated `upstream.lock` after first accepted build.

Steps:

- [x] Add flavor manifests for the full milestone 1 matrix.
- [x] Add lean Dockerfiles that install only matching engine/vendor dependencies and verify the runtime command during image construction.
- [x] Declare Ray as an upstream only for engine/vendor modes that need it.
- [x] Add labels for engine, vendor, flavor, source, revision, and build date.
- [x] Ensure no disallowed `*-ray`, `all-engines`, or `all-gpu-vendors` flavor exists.

## Task 4: image build and publish workflows

Files:

- `.github/workflows/engine-detect.yml`.
- `.github/workflows/engine-build.yml`.
- `.github/workflows/engine-publish.yml`.
- `ci/build/build-image.sh`.
- `ci/build/publish-image.sh`.

Interfaces:

- PR workflows build and test but do not publish stable tags.
- Manual/scheduled release workflows publish by digest and attach flavor tags only after CT passes.

Steps:

- [x] Run the detector and expose its matrix to downstream jobs.
- [x] Build each selected flavor with a flavor-scoped cache.
- [x] Push build artifacts by digest for CT without moving `latest`.
- [x] Attach immutable and `*-latest` tags only in publish workflow after gates pass.
- [x] Upload SBOM, provenance, detector plan, and CT report artifacts.

## Task 5: static and vendor container tests

Files:

- `ci/container-tests/static-check.sh`.
- `ci/container-tests/openai-chat-smoke.sh`.
- `ci/container-tests/vendor-device-check.sh`.
- `internal/ct/report.go`.
- `internal/ct/report_test.go`.
- `docs/container-test-strategy.md`.

Interfaces:

- Test report schema `kuvryn.container-test-report/v1`.

Steps:

- [x] Implement static image checks for labels, entrypoint, expected commands, wrong vendor stacks, and SBOM generation.
- [x] Implement vendor device checks for NVIDIA, AMD, Intel, and CPU.
- [x] Implement one real inference smoke where the engine supports OpenAI-compatible APIs.
- [x] Emit normalized JSON reports.
- [x] Make hardware-certified publishing fail closed when required vendor CT is missing or failed; mark AMD/Intel outputs hardware-untested while runners are unavailable.

## Task 6: distributed CT and collective validation

Files:

- `ci/container-tests/distributed-smoke.sh`.
- `internal/ct/distributed.go`.
- `internal/ct/distributed_test.go`.
- `docs/container-test-strategy.md`.

Interfaces:

- Distributed CT consumes flavor capabilities and validates ranks, world size, collective backend, and one inference response.

Steps:

- [x] Define distributed CT fixtures for two-rank minimum tests.
- [x] Wire NCCL validation for NVIDIA distributed-capable flavors on NVIDIA runners.
- [x] Represent RCCL for AMD distributed-capable flavors as build-only/hardware-untested because no AMD hardware is available.
- [x] Represent HCCL/oneCCL path for Intel distributed-capable flavors as build-only/hardware-untested because no Intel hardware is available.
- [x] Record limitations when distributed support is not yet proven for a vendor/engine combination.

## Task 7: cluster network and RDMA discovery

Files:

- `internal/clusterprobe/probe.go`.
- `internal/clusterprobe/parser.go`.
- `internal/clusterprobe/topology.go`.
- `internal/clusterprobe/rdma.go`.
- `internal/clusterprobe/*_test.go`.
- `docs/cluster-network-rdma-design.md`.

Interfaces:

- Host capability record schema covering accelerator, network, fabric, runtime, limitations, and collectives.

Steps:

- [x] Implement sentinel-delimited probe parsing.
- [x] Implement read-only SSH/direct probe script generation.
- [x] Implement Kubernetes diagnostic Job/DaemonSet design stubs and parser surfaces.
- [x] Classify single-node, direct, ring, switch, Kubernetes fabric, ethernet-only, and unknown.
- [x] Implement RDMA validation result types for perftest and collective bootstrap outcomes.
- [x] Test malformed, partial, and complete probe output.

## Task 8: collective environment and placement constraints

Files:

- `internal/placement/constraints.go`.
- `internal/placement/collectives.go`.
- `internal/placement/compatibility.go`.
- `internal/placement/*_test.go`.

Interfaces:

- `func PlanCollectiveEnv(record ClusterRecord, flavor ImageManifest) (CollectiveEnvPlan, error)`.
- `func CheckPlacement(req DeploymentRequest, inventory ClusterRecord) (PlacementPlan, error)`.

Steps:

- [x] Generate NCCL env from validated NVIDIA fabric facts.
- [x] Generate RCCL env from validated AMD fabric facts.
- [x] Generate HCCL/oneCCL env from validated Intel fabric facts.
- [x] Fail RDMA-required placements when validation is missing, failed, stale, or pod-invisible.
- [x] Refuse mixed-vendor auto-placement without explicit layout.

## Task 9: runtime wrapper and operability contract

Files:

- `cmd/kuvryn-runtime-wrapper/main.go`.
- `internal/runtimecontract/state.go`.
- `internal/runtimecontract/health.go`.
- `internal/runtimecontract/metrics.go`.
- `internal/runtimecontract/*_test.go`.
- `docs/runtime-operability-tooling.md`.

Interfaces:

- `/kuvryn/health/live`, `/kuvryn/health/ready`, `/kuvryn/health/startup`.
- `/metrics` Prometheus endpoint or adapter.
- `/var/run/kuvryn/runtime-state.json`.
- `/var/log/kuvryn/events.ndjson`.

Steps:

- [x] Build a tiny wrapper that starts one runtime process, forwards signals, reaps children, and exits with the runtime status.
- [x] Write structured lifecycle events.
- [x] Expose live/startup/ready checks with engine adapter hooks.
- [x] Add normalized metrics baseline.
- [x] Test graceful shutdown, child reaping, failed startup, and readiness not treating port-open as inference-ready.

## Task 10: OOM, heartbeat, and debug evidence integration

Files:

- `internal/ops/heartbeat.go`.
- `internal/ops/oom.go`.
- `internal/ops/debug_bundle.go`.
- `internal/ops/*_test.go`.
- `docs/runtime-operability-tooling.md`.

Interfaces:

- Heartbeat payload includes deployment id, image digest, flavor, node/pod/container identity, state, readiness result, and rank/world data.
- Debug bundle is opt-in and secret-redacted.

Steps:

- [x] Define heartbeat payload and state transitions.
- [x] Parse Docker/local OOM evidence from container exit, cgroup memory events, and daemon inspect data.
- [x] Parse Kubernetes OOM/eviction evidence from pod status and events.
- [x] Redact secrets from debug bundles.
- [x] Test runtime crash, cgroup OOM, host OOM limitation, Kubernetes eviction, health restart, and user stop classification.

## Task 11: documentation, examples, and release readiness

Files:

- `README.md`.
- `docs/*.md`.
- `.procoder/backlog/**`.
- `.procoder/specs/engine-runtime-foundation.md`.
- `.procoder/plans/engine-runtime-foundation.md`.

Interfaces:

- User-facing quick start for adding a flavor, detecting updates, building, testing, and publishing.

Steps:

- [x] Add README overview and command examples.
- [x] Add example manifests and locks for one CPU and one GPU flavor.
- [x] Cross-link design docs from the README.
- [x] Keep backlog stories aligned with implemented scope.
- [x] Run `procoder test`, `procoder review`, and `procoder check` before calling implementation complete.

## Completion evidence

- Repository implementation covers the planned schema, detector, manifest, Dockerfile, CI/CD/CT, distributed CT, probe, placement, runtime wrapper, OOM/debug, docs, and example lock slices.
- `bash -n` passed for touched shell scripts including `ci/container-tests/distributed-smoke.sh`, runtime smoke scripts, and image install/verify helpers.
- `go test ./...` passed.
- `go run ./cmd/kuvryn-manifest-check --root .` validated all engine manifests.
- `procoder test`, `procoder lint`, `procoder security`, and `procoder check` passed with no blockers.
- NVIDIA DGX CT evidence is stored in `ct-reports/`; AMD/Intel hardware CT is intentionally out of local scope and represented as build-only `hardware-untested`.
