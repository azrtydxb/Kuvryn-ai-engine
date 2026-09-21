# Runtime tooling baseline is present without heavy platforms

Status: done
Created: 2026-09-16
Epic: runtime-operability-contract
Sprint: -

## Description

Ensure every engine image contains the agreed small runtime-local toolkit and excludes heavy operations platforms.

## Acceptance criteria

- [x] Baseline tools include `tini`, `curl`, `jq`, `procps`, `iproute2`, and `ca-certificates` unless a base image already provides a verified equivalent.
- [x] Kuvryn runtime wrapper, health/readiness adapter, metrics adapter, lifecycle logging, and version/build manifest are present.
- [x] Full systemd, Prometheus server, Grafana, Alertmanager, SSH server, databases, queues, and unrelated GPU vendor stacks are absent.
- [x] Static CT enforces required and excluded tool lists.

## Evidence

## Closure evidence

- `go test ./...` passed for schema, detector, CT, placement, runtime wrapper, ops, and publish tag packages.
- `go run ./cmd/kuvryn-manifest-check --root .` validated all 11 engine manifests.
- `procoder test`, `procoder lint`, `procoder security`, and `procoder check` passed with 0 blockers.
- DGX Spark NVIDIA CT evidence was captured under `ct-reports/` for NVIDIA hardware-certified flavors; AMD/Intel remain explicitly build-only and `hardware-untested` because matching hardware is unavailable.
- Accepted upstream locks under `engines/**/upstream.lock` record published digests, CT report paths, and vendor-before-engine publish tags for CPU/NVIDIA certified flavors.
