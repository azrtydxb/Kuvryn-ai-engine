# Metrics heartbeat and OOM evidence are normalized

Status: done
Created: 2026-09-16
Epic: runtime-operability-contract
Sprint: -

## Description

Define and implement normalized runtime metrics, heartbeat payloads, and OOM/failure evidence classification for Docker/local and Kubernetes deployments.

## Acceptance criteria

- [x] Metrics include uptime, request counts, errors, active requests, tokens, TTFT/latency, model load duration, memory, restart/OOM count, rank health, and fabric status where available.
- [x] Heartbeat includes deployment id, image digest, flavor, node/pod/container identity, state, readiness result, and rank/world info.
- [x] Docker/local OOM evidence uses container exit, cgroup memory events, daemon inspect, and kernel/journal where available.
- [x] Kubernetes evidence uses pod status, termination reason, events, and kubelet messages where available.
- [x] Debug bundles redact secrets and avoid environment dumps.

## Evidence

## Closure evidence

- `go test ./...` passed for schema, detector, CT, placement, runtime wrapper, ops, and publish tag packages.
- `go run ./cmd/kuvryn-manifest-check --root .` validated all 11 engine manifests.
- `procoder test`, `procoder lint`, `procoder security`, and `procoder check` passed with 0 blockers.
- DGX Spark NVIDIA CT evidence was captured under `ct-reports/` for NVIDIA hardware-certified flavors; AMD/Intel remain explicitly build-only and `hardware-untested` because matching hardware is unavailable.
- Accepted upstream locks under `engines/**/upstream.lock` record published digests, CT report paths, and vendor-before-engine publish tags for CPU/NVIDIA certified flavors.
