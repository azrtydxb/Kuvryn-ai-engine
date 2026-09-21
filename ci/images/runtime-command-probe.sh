#!/usr/bin/env bash
set -euo pipefail

flavor=${1:?usage: runtime-command-probe.sh <flavor>}
case "$flavor" in
vllm-*) command -v vllm >/dev/null && python3 -c 'import vllm' ;;
sglang-*) command -v sglang >/dev/null && python3 -c 'import sglang' ;;
llama-cpp-*) command -v llama-server >/dev/null && llama-server --help >/dev/null ;;
tensorrt-llm-*) command -v trtllm-serve >/dev/null && { trtllm-serve --help >/dev/null || trtllm-serve --version >/dev/null; } ;;
*)
	echo "unknown flavor: ${flavor}" >&2
	exit 2
	;;
esac
