# Task 4: image build and publish workflows

Status: done
Created: 2026-09-16
Plan: .procoder/plans/engine-runtime-foundation.md

## Description

Wire GitHub Actions and scripts for detector, build, CT evidence, and gated publish.

## Acceptance criteria

- [x] Detector workflow emits the update plan artifact.
- [x] Build workflow consumes detector output to select changed flavors.
- [x] Cache scopes are flavor-specific.
- [x] PR workflows do not publish stable tags.
- [x] Publish workflow attaches tags only after CT gates pass.

## Evidence

- `.github/workflows/engine-detect.yml` uploads `update-plan.json`.
- `.github/workflows/engine-build.yml` derives `changed-flavors.txt` from detector output, can also build one manual flavor, only runs after a successful non-PR detect workflow, and checks out the triggering workflow SHA.
- `ci/build/build-image.sh` uses `.buildx-cache/<flavor>` local buildx cache scopes where the buildx driver supports local cache export, falls back safely for Docker-driver DGX hosts, and when `BUILD_PUSH=1`, pushes a temporary `ci-<sha>-<flavor>` tag while recording `image-digests/<flavor>.digest` for downstream CT/publish.
- `ci/build/publish-image.sh` runs manifest and CT validation before `docker buildx imagetools create` tag movement.
- `cmd/kuvryn-publish-tags` generates digest, `v<kuvryn-version>_<vendor>_<engine>_<engine-version>`, and `v<kuvryn-version>_<vendor>_<engine>_nightly` tags from `VERSION`, manifest vendor, and the manifest engine upstream.
- GitHub Actions are pinned to commit SHAs.
