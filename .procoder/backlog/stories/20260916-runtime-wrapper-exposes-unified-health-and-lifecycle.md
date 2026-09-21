# Runtime wrapper exposes unified health and lifecycle

Status: done
Created: 2026-09-16
Epic: runtime-operability-contract
Sprint: -

## Description

Implement the minimal runtime wrapper that starts one engine process, handles signals, writes lifecycle events, exposes health/readiness/startup endpoints, and avoids full systemd.

## Acceptance criteria

- [x] Wrapper starts the configured engine command and exits with its status.
- [x] Wrapper forwards SIGTERM/SIGINT and enforces graceful shutdown timeout.
- [x] Wrapper reaps child processes.
- [x] `/kuvryn/health/live`, `/kuvryn/health/ready`, and `/kuvryn/health/startup` exist or are proxied through an adapter.
- [x] Lifecycle state and events are written to the documented files.
- [x] Tests cover startup failure, graceful shutdown, child reaping, and readiness not equaling port-open.

## Evidence

## Closure evidence

- `go test ./...` passed for schema, detector, CT, placement, runtime wrapper, ops, and publish tag packages.
- `go run ./cmd/kuvryn-manifest-check --root .` validated all 11 engine manifests.
- `procoder test`, `procoder lint`, `procoder security`, and `procoder check` passed with 0 blockers.
- DGX Spark NVIDIA CT evidence was captured under `ct-reports/` for NVIDIA hardware-certified flavors; AMD/Intel remain explicitly build-only and `hardware-untested` because matching hardware is unavailable.
- Accepted upstream locks under `engines/**/upstream.lock` record published digests, CT report paths, and vendor-before-engine publish tags for CPU/NVIDIA certified flavors.
