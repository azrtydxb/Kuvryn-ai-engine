# Task 1: repository foundation and schemas

Status: done
Created: 2026-09-16
Plan: .procoder/plans/engine-runtime-foundation.md

## Description

Bootstrap the Go module and implement strict manifest/lock schema loading and validation.

## Acceptance criteria

- [x] Go module is initialized for `github.com/azrtydxb/kuvryn-ai-engine`.
- [x] `internal/engineimage` defines manifest and lock structs.
- [x] YAML loading rejects unknown fields.
- [x] Validation covers flavor naming, engine/vendor compatibility, required tests, and GHCR namespace.
- [x] Unit tests cover valid and invalid schema cases.

## Evidence

- `go test ./...` passes.
- `internal/engineimage/schema.go` and `internal/engineimage/validate.go` implement manifest/lock loading and validation, including accepted image digest, publish tags, and CT evidence fields.
- `internal/engineimage/schema_test.go` covers valid manifest loading, unknown YAML fields, invalid schema cases, malformed locks, and example locks.
