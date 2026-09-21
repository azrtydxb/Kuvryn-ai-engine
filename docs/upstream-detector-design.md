# Upstream detector design

## Purpose

The upstream detector decides which engine image flavors need a rebuild. It is not a builder and it is not a publisher. Its only output is a deterministic plan: changed inputs, affected flavors, and a CI matrix.

The detector is implemented in Go first. Python is an allowed fallback only for practical API or registry integration gaps.

## Inputs

The detector reads:

```text
engines/**/image.yaml
engines/**/upstream.lock
.github/workflows/**
ci/**
```

Each `image.yaml` declares the upstreams and local files that affect a flavor. Each `upstream.lock` records the last accepted resolved inputs.

## Output files

The detector writes machine-readable JSON to stdout or an explicit output path:

```json
{
  "schema": "kuvryn.engine-update-plan/v1",
  "changed": true,
  "flavors": [
    {
      "name": "vllm-nvidia",
      "engine": "vllm",
      "vendor": "nvidia",
      "path": "engines/vllm/nvidia",
      "runner": "self-hosted,nvidia,gpu",
      "reason": [
        "upstream:vllm",
        "base-image:nvcr.io/nvidia/pytorch:25.01-py3"
      ],
      "publish": false
    }
  ],
  "matrix": {
    "include": [
      {
        "name": "vllm-nvidia",
        "engine": "vllm",
        "vendor": "nvidia",
        "path": "engines/vllm/nvidia",
        "runner": "self-hosted,nvidia,gpu"
      }
    ]
  }
}
```

`publish` defaults to false in pull requests. Release, scheduled, or manual workflows may enable publishing only after build and CT pass.

## Change detection rules

A flavor is marked affected when any of these change:

| Input               | Detection method                         | Scope                                              |
| ------------------- | ---------------------------------------- | -------------------------------------------------- |
| Engine release      | GitHub release/tag/commit comparison     | One engine/vendor flavor unless shared by variants |
| Base image          | Registry digest comparison               | Flavors using that exact base image reference      |
| Accelerator stack   | Manifest-declared version or base digest | Vendor-specific flavors only                       |
| Dockerfile          | Local file digest                        | Owning flavor                                      |
| Shared build script | Local file digest declared in manifest   | Declared consuming flavors                         |
| CT script           | Local file digest declared in manifest   | Declared consuming flavors                         |
| Manifest            | Local file digest                        | Owning flavor                                      |
| Lock schema version | Schema comparison                        | Owning flavor                                      |

A flavor is not affected by another flavor's upstream unless it declares that upstream as an input.

## Supported upstream resolvers

Milestone 1 should support these resolvers:

| Type               | Required fields                  | Resolved identity                       |
| ------------------ | -------------------------------- | --------------------------------------- |
| `github-release`   | `repo`, optional `constraint`    | tag, commit, published time             |
| `github-tag`       | `repo`, `tag` or `constraint`    | tag, commit                             |
| `github-commit`    | `repo`, `ref`                    | commit SHA                              |
| `container-digest` | `image`, `tag`                   | digest, media type, platform digests    |
| `pypi-version`     | `package`, optional `constraint` | version, artifact hashes when available |
| `static-version`   | `version`                        | literal value                           |
| `local-file`       | `path`                           | SHA-256 digest                          |

Unsupported resolver types fail closed. They do not silently mark a flavor unchanged.

## Manifest integration

Example `image.yaml` upstream block:

```yaml
upstreams:
  - name: vllm
    type: github-release
    repo: vllm-project/vllm
    constraint: stable
  - name: base
    type: container-digest
    image: nvcr.io/nvidia/pytorch
    tag: 25.01-py3
  - name: dockerfile
    type: local-file
    path: engines/vllm/nvidia/Dockerfile
  - name: smoke-test
    type: local-file
    path: ci/container-tests/openai-chat-smoke.sh
```

Resolved values are compared against `upstream.lock`.

## Lock update contract

The detector may produce a proposed lock update, but it must not commit it by itself. `kuvryn-lock-accept` writes the accepted `upstream.lock` after publish evidence has passed validation.

Lock files are updated only after:

1. the flavor builds successfully;
2. static image checks pass;
3. `kuvryn-ct-validate` accepts the CT report against the flavor manifest;
4. matching vendor CT passes, or the flavor is explicitly marked as CT-deferred by policy;
5. the published image digest is known.

Current CT-deferred policy: AMD and Intel flavors are built even though local matching GPU runners are not available yet. They use `publish.certification: hardware-untested`, are reported as hardware-untested, and must not be promoted as hardware-certified until their matching vendor CT passes. Their lock evidence may record static CT acceptance without vendor CT while this policy is active.

This prevents recording untested upstreams as accepted state.

## CI integration

### Pull request

```text
detect -> build affected flavors -> static CT -> vendor CT where runners are available -> report
```

Pull requests do not publish stable tags.

### Scheduled upstream check

```text
detect upstream drift -> open/update PR with changed locks and metadata
```

The scheduled workflow should prefer opening a PR over pushing directly. The PR then uses the same build and CT gates.

### Manual release

```text
detect -> build -> CT -> publish by digest -> attach flavor tags -> update lock files
```

Manual release is the path that may publish `*-latest` tags.

## Runner mapping

| Vendor | Runner label                       | Required checks                                                             |
| ------ | ---------------------------------- | --------------------------------------------------------------------------- |
| NVIDIA | `self-hosted,nvidia,gpu`           | CUDA visibility, NCCL when distributed, inference smoke                     |
| AMD    | `self-hosted,amd,gpu`              | ROCm visibility, RCCL when distributed, inference smoke                     |
| Intel  | `self-hosted,intel,gpu`            | Intel accelerator visibility, HCCL/oneCCL when distributed, inference smoke |
| CPU    | `ubuntu-latest` or self-hosted CPU | CPU inference smoke                                                         |

GPU CT fails closed for hardware-certified release publishing. If a runner is missing, the flavor can still be built and published only under the explicit CT-deferred, hardware-untested policy; it must not be silently treated as validated and must not move a `latest` tag.

## Failure modes

| Failure                     | Behavior                                                                       |
| --------------------------- | ------------------------------------------------------------------------------ |
| Upstream API unavailable    | Fail detection unless workflow explicitly allows offline lock-only comparison. |
| Registry digest unavailable | Fail detection for affected flavor.                                            |
| Unsupported upstream type   | Fail detection.                                                                |
| Missing lock file           | Mark flavor changed and require full build/CT.                                 |
| Malformed manifest          | Fail before matrix generation.                                                 |
| Changed flavor has no tests | Fail plan generation.                                                          |

## Security and reproducibility

- Do not print tokens or registry credentials.
- Prefer immutable digests over mutable tags in locks.
- Publish SBOM and provenance per flavor.
- Do not execute commands from manifests during detection.
- Validate paths stay inside the repository.
- Treat untrusted external metadata as data, not shell input.
