# Task 7: cluster network and RDMA discovery

Status: done
Created: 2026-09-16
Plan: .procoder/plans/engine-runtime-foundation.md

## Description

Implement probe parsing, topology classification, and RDMA validation result scaffolding.

## Acceptance criteria

- [x] Sentinel-delimited probe parsing is implemented.
- [x] SSH/direct read-only probe generation exists.
- [x] Kubernetes diagnostic surfaces are represented.
- [x] Topology classifier covers all required classes.
- [x] Tests cover malformed, partial, and complete probe output.

## Evidence

- `internal/clusterprobe/parser.go` parses sentinel-delimited sections.
- `internal/clusterprobe/probe.go` exposes direct probe script generation and `KubernetesDiagnosticCommand`.
- `internal/clusterprobe/topology.go` and `internal/clusterprobe/rdma.go` represent topology/RDMA outcomes.
- `internal/clusterprobe/clusterprobe_test.go` covers malformed parsing, complete sections, Kubernetes diagnostic command, topology, and pod-visible RDMA behavior.
