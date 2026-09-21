# Runtime operability tooling

## Purpose

Every engine image and deployment should expose health, readiness, metrics, heartbeat, and failure evidence in a unified way without turning engine images into large all-in-one platforms.

The principle is: keep engine images lean, but require a small common operability contract around every runtime.

## Recommended components

| Component                | Where it should live                         | Required         | Purpose                                                                                      |
| ------------------------ | -------------------------------------------- | ---------------- | -------------------------------------------------------------------------------------------- |
| Runtime wrapper          | Inside each engine image                     | Yes              | Starts the engine, normalizes signals, writes lifecycle state, handles graceful shutdown.    |
| Health/readiness adapter | Inside each engine image                     | Yes              | Exposes a common health contract even when engines use different endpoints.                  |
| Metrics adapter          | Sidecar or embedded lightweight exporter     | Yes              | Normalizes engine, GPU, memory, queue, token, and request metrics.                           |
| Heartbeat                | Control-plane agent or sidecar               | Yes              | Proves the process/pod/node is alive and still owned by the current deployment.              |
| OOM detector             | Node agent plus container/pod event reader   | Yes              | Captures cgroup/kernel/Kubernetes OOM evidence and links it to the deployment.               |
| Log/event collector      | Sidecar, node agent, or platform integration | Yes              | Collects structured startup, readiness, error, and shutdown events.                          |
| Supervisor               | Container entrypoint or orchestrator         | Yes, but minimal | Handles one main runtime process and shutdown; do not run full systemd in normal containers. |
| Debug bundle collector   | Control-plane triggered                      | Yes              | Captures safe diagnostics after failure without dumping secrets.                             |
| OpenTelemetry exporter   | Sidecar or direct integration                | Recommended      | Standard traces/metrics/logs integration.                                                    |
| Prometheus exporter      | Sidecar or direct integration                | Recommended      | Common scraping path for Kubernetes and bare-metal clusters.                                 |

## What belongs in the image

Each engine image should include only the minimum runtime-local tooling:

- a small entrypoint/wrapper;
- health/readiness adapter;
- metrics endpoint or exporter hook;
- signal handling and graceful shutdown;
- version/manifest endpoint or file;
- optional engine-specific probes required by that engine.

The image should not include a full monitoring stack, dashboard, database, alertmanager, or unrelated vendor tooling.

## What belongs outside the image

These should usually be sidecars, DaemonSets, host agents, or control-plane services:

- GPU/node exporters;
- cAdvisor/container metrics;
- Kubernetes event watchers;
- OOM/kmsg/journal collectors;
- log shipping;
- Prometheus, Grafana, Alertmanager;
- long-running remediation agents;
- cluster-level heartbeat aggregation.

This keeps engine images small and lets operations tooling update independently from vLLM, SGLang, TensorRT-LLM, and llama.cpp.

## Unified runtime contract

Every runtime deployment must expose or produce:

```yaml
runtime_contract:
  health:
    live: /kuvryn/health/live
    ready: /kuvryn/health/ready
    startup: /kuvryn/health/startup
  metrics:
    prometheus: /metrics
  lifecycle:
    state_file: /var/run/kuvryn/runtime-state.json
    events_file: /var/log/kuvryn/events.ndjson
  identity:
    image_digest: required
    engine: required
    vendor: required
    model: required
    deployment_id: required
```

The adapter may proxy engine-native endpoints, for example vLLM or SGLang `/health`, but the control plane should consume the Kuvryn contract.

## Health model

| Check    | Meaning                                                   | Must verify                                                            |
| -------- | --------------------------------------------------------- | ---------------------------------------------------------------------- |
| Startup  | Runtime process entered startup and has not failed early. | Process exists, config parsed, model load started.                     |
| Live     | Runtime process should be restarted if false.             | Main process alive, event loop responsive.                             |
| Ready    | Runtime can serve traffic.                                | Engine health plus real inference readiness where applicable.          |
| Degraded | Runtime can serve but has a problem.                      | Optional: missing replica, slow fabric, degraded GPU, high error rate. |

A TCP port alone is not ready. For LLM serving, readiness should include an inference probe when the engine supports it.

## Metrics baseline

All runtimes should normalize at least:

- process uptime;
- request count, error count, active requests;
- queue depth or waiting requests where available;
- tokens generated, prompt tokens, output tokens where available;
- time to first token and end-to-end latency;
- model load duration;
- GPU memory used/free per device;
- CPU memory and cgroup memory;
- OOM/restart count;
- distributed rank health;
- collective/fabric validation status.

Engine-native metrics should be preserved, but dashboards and alerts should use the normalized Kuvryn names where possible.

## OOM and failure evidence

OOM handling must capture evidence from the layer that actually observed the kill:

| Environment  | Evidence source                                                                          |
| ------------ | ---------------------------------------------------------------------------------------- |
| Docker/local | container exit code, cgroup memory events, daemon inspect, kernel/journal when available |
| Kubernetes   | pod status, container terminated reason, node events, kubelet messages when available    |
| Host agent   | dmesg/journal/cgroup files, with secret-safe filtering                                   |

The control plane should distinguish:

- runtime crash;
- CUDA/ROCm/driver crash;
- cgroup OOM;
- host OOM;
- Kubernetes eviction;
- health-check restart;
- user-requested stop.

## Supervisor policy

Use a tiny supervisor behavior, not full systemd, inside normal containers:

- run one main runtime process;
- forward SIGTERM/SIGINT;
- enforce graceful shutdown timeout;
- reap child processes;
- write structured lifecycle events;
- exit with the runtime's failure code.

Full systemd-like behavior belongs at the node or orchestrator layer, not inside every engine image, unless a specific engine explicitly requires it.

## Heartbeat model

Heartbeats should include:

- deployment id;
- image digest;
- engine/vendor/flavor;
- node/pod/container identity;
- runtime pid or container id;
- current state: starting, ready, degraded, draining, stopped, failed;
- last readiness result;
- rank/world information for distributed workloads.

A heartbeat should never be treated as readiness. It only proves the observer path is alive.

## Security rules

- Do not expose environment dumps.
- Redact tokens, API keys, model registry credentials, and cluster credentials.
- Keep debug bundles opt-in and scoped to a deployment.
- Bind internal health/metrics endpoints to localhost or pod network unless explicitly exposed.
- Separate user-facing inference auth from internal metrics/health auth.

## Implementation guidance

Milestone 1 should implement the contract shape and static checks first:

1. define the runtime state JSON schema;
2. define normalized metric names;
3. add static image checks that verify labels, entrypoint, health contract, and no wrong vendor stack;
4. add a simple wrapper pattern for each engine family;
5. add control-plane ingestion later.

This gives every engine a consistent operational surface without blocking the image build foundation on a full observability platform.
