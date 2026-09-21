# Container test strategy

## Purpose

Container tests prove that each image flavor is usable on its intended hardware without relying on an all-in-one image. Tests are flavor-scoped and vendor-aware.

## Test levels

| Level              | Hardware required       | Blocks PR                    | Blocks publish                  | Purpose                                                                  |
| ------------------ | ----------------------- | ---------------------------- | ------------------------------- | ------------------------------------------------------------------------ |
| Static image check | No                      | Yes                          | Yes                             | Validate metadata, commands, labels, and absence of wrong vendor stacks. |
| CPU smoke          | CPU only                | Yes for CPU flavors          | Yes for CPU flavors             | Run tiny llama.cpp or endpoint smoke where possible.                     |
| Vendor smoke       | Matching GPU            | Yes when runner is available | Yes for hardware-certified tags | Start server and complete one inference request.                         |
| Distributed CT     | Matching multi-node GPU | For distributed changes      | Yes for hardware-certified tags | Validate sharding, collectives, and rank startup.                        |
| Benchmark sanity   | Matching GPU            | No, report initially         | Optional release gate later     | Catch extreme regressions before publication.                            |

## Static image checks

Every flavor must pass static checks before any GPU time is consumed.

Required checks:

- image has OCI labels for engine, vendor, version, source, revision, and build date;
- expected runtime command exists and is executable (`runtime-command` is required evidence, not informational);
- expected Python package or binary imports/starts;
- image does not contain unrelated accelerator stacks;
- image exposes or documents the expected serving port;
- image has no known test secrets or credentials in environment defaults;
- SBOM generation succeeds.

Examples of wrong-stack failures:

| Flavor          | Must not contain                                                          |
| --------------- | ------------------------------------------------------------------------- |
| `vllm-nvidia`   | ROCm tools, oneAPI runtime, SGLang package                                |
| `vllm-amd`      | CUDA toolkit, NVIDIA driver libraries copied into image, oneAPI runtime   |
| `sglang-intel`  | CUDA toolkit, ROCm tools, vLLM package                                    |
| `llama-cpp-cpu` | CUDA toolkit, ROCm tools, oneAPI runtime unless explicitly variant-scoped |

## Vendor smoke checks

Current runner policy: NVIDIA is the only hardware-backed CT target because this project does not currently have AMD or Intel hardware. AMD and Intel images are still built with `publish.certification: hardware-untested`; their `vendor-device` checks remain declared in manifests with `required: false` so the CI shape is visible without pretending validation happened.

Hardware-certified GPU flavors run on a matching hardware runner. Today that means NVIDIA only:

| Vendor | Runner label             | Expected device check                        | Current status                   |
| ------ | ------------------------ | -------------------------------------------- | -------------------------------- |
| NVIDIA | `self-hosted,nvidia,gpu` | `nvidia-smi` and CUDA-visible device         | hardware CT path configured      |
| AMD    | not configured           | `rocminfo` on future AMD hardware            | build-only / `hardware-untested` |
| Intel  | not configured           | Intel accelerator tool on future Intel hosts | build-only / `hardware-untested` |

Smoke flow:

```text
pull built image by digest
start runtime with tiny model or synthetic fixture
wait for runtime health endpoint
send one real inference request
record latency and time-to-first-token when streaming exists
stop container and collect logs on failure
```

Use `ci/container-tests/prepare-smoke-fixture.sh <flavor> <env-file>` to prepare the default tiny fixture, then pass that env file to `ci/container-tests/engine-runtime-smoke.sh` with `KUVRYN_SMOKE_ENV_FILE=<env-file>`. The default fixtures are `facebook/opt-125m` for vLLM/SGLang and `aladar/tiny-random-LlamaForCausalLM-GGUF` for llama.cpp. Override with `KUVRYN_HF_TINY_MODEL` when validating a specific HF fixture. Smoke-only runtime flags can be appended with `KUVRYN_VLLM_SERVE_ARGS`, `KUVRYN_SGLANG_SERVE_ARGS`, or `KUVRYN_LLAMA_CPP_SERVE_ARGS` for constrained hardware. Set `KUVRYN_TRTLLM_SERVE_ARGS` for TensorRT-LLM because it serves a built engine/checkpoint rather than a raw HF model.

A health endpoint alone is not enough. A publishable image must produce tokens or the runtime-specific equivalent response.

## Distributed CT

Distributed CT runs only for flavors declaring distributed support.

Required distributed checks:

- two or more ranks start successfully;
- expected collective backend initializes: NCCL, RCCL, HCCL, or oneCCL;
- ranks agree on world size and rank identity;
- tensor-parallel or pipeline-parallel mode serves one inference request;
- failure logs include rank, host, engine, vendor, and image digest.

Distributed CT may use Kubernetes, Docker multi-host, or a backend-specific launcher, but the test result must be normalized into the same report shape.

## Test report shape

Each test writes JSON so CI can summarize failures without scraping logs. `kuvryn-ct-validate` checks this evidence against the flavor manifest before publish tagging.

```json
{
  "schema": "kuvryn.container-test-report/v1",
  "flavor": "vllm-nvidia",
  "image": "ghcr.io/azrtydxb/kuvryn-ai-engine/vllm-nvidia@sha256:example",
  "vendor": "nvidia",
  "tests": [
    {
      "name": "openai-chat-smoke",
      "status": "passed",
      "duration_ms": 12345,
      "metrics": {
        "ttft_ms": 800
      }
    }
  ]
}
```

## Publish gate

An image can be published as hardware-certified only when:

1. static checks pass;
2. required vendor smoke checks pass;
3. required distributed CT passes for distributed-capable tags;
4. SBOM and provenance are written under `publish-artifacts/<flavor>/` and uploaded with the accepted lock;
5. the image digest and allowed publish tags are recorded in `upstream.lock`.

AMD and Intel flavors may be built and published only as hardware-untested artifacts in this repository because no AMD/Intel hardware is available. The publish script runs `kuvryn-ct-validate`, tags those artifacts with `hardware-untested`, and does not move `latest`. Release notes, tags, and lock/test evidence must not describe them as hardware-certified. Missing hardware is not a hardware pass; it is an explicit deferred-validation state.
