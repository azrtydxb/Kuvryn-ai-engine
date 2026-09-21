#!/usr/bin/env bash
set -euo pipefail

if ! command -v trtllm-serve >/dev/null; then
	echo "trtllm-serve is not present in the TensorRT-LLM base image" >&2
	exit 1
fi
if [[ "${TENSORRT_LLM_SKIP_HELP:-0}" != "1" ]]; then
	trtllm-serve --help >/dev/null || trtllm-serve --version >/dev/null
fi
