# Task 8: collective env and placement

Status: done
Created: 2026-09-16
Plan: .procoder/plans/engine-runtime-foundation.md

## Description

Add collective backend environment planning and fail-closed placement checks.

## Acceptance criteria

- [x] NCCL, RCCL, HCCL/oneCCL, and no-RDMA paths are represented.
- [x] RDMA-required placements fail when validation is missing or failed.
- [x] Kubernetes pod-invisible RDMA fails RDMA-required placement.
- [x] Mixed-vendor auto-placement fails without explicit layout.
- [x] Tests cover vendor and failure paths.

## Evidence

- `internal/placement/collectives.go` handles `none`, `nccl`, `rccl`, `hccl`, and `oneccl`.
- `internal/placement/constraints.go` rejects mixed-vendor automatic placement and RDMA-required placement without fresh usable validation.
- `internal/placement/placement_test.go` covers NCCL env, mixed vendor rejection, missing RDMA validation, successful RDMA validation, and Kubernetes pod-invisible RDMA rejection.
