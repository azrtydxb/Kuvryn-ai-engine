#!/usr/bin/env bash
set -euo pipefail

flavor=${1:?usage: prepare-smoke-fixture.sh <flavor> [output-env-file]}
out=${2:-/tmp/kuvryn-smoke-fixture.env}
cache=${KUVRYN_MODEL_CACHE:-$PWD/.kuvryn-model-cache}
mkdir -p "$cache"

write_env() {
	: >"$out"
	for kv in "$@"; do
		printf '%s\n' "$kv" >>"$out"
	done
	printf '%s\n' "$out"
}

case "$flavor" in
llama-cpp-*)
	repo=${KUVRYN_LLAMA_GGUF_REPO:-aladar/tiny-random-LlamaForCausalLM-GGUF}
	file=${KUVRYN_LLAMA_GGUF_FILE:-tiny-random-LlamaForCausalLM.gguf}
	model_dir="$cache/llama-cpp/${repo//\//__}"
	model_path="$model_dir/$file"
	mkdir -p "$model_dir"
	if [[ ! -s "$model_path" ]]; then
		curl -fL --retry 3 --retry-delay 2 -o "$model_path" "https://huggingface.co/${repo}/resolve/main/${file}"
	fi
	write_env "KUVRYN_SMOKE_MODEL=/models/$file" "KUVRYN_SMOKE_MODEL_NAME=$file" "KUVRYN_MODEL=/models/$file" "KUVRYN_SMOKE_MODEL_MOUNT=$model_dir:/models:ro"
	;;
vllm-* | sglang-*)
	model=${KUVRYN_HF_TINY_MODEL:-facebook/opt-125m}
	write_env "KUVRYN_SMOKE_MODEL=$model" "KUVRYN_SMOKE_MODEL_NAME=$model" "KUVRYN_MODEL=$model"
	;;
tensorrt-llm-*)
	if [[ -z "${KUVRYN_TRTLLM_SERVE_ARGS:-}" ]]; then
		echo "TensorRT-LLM smoke needs a prebuilt tiny TensorRT-LLM engine/checkpoint." >&2
		echo "Set KUVRYN_TRTLLM_SERVE_ARGS to the trtllm-serve arguments for that fixture." >&2
		exit 2
	fi
	write_env "KUVRYN_TRTLLM_SERVE_ARGS=$KUVRYN_TRTLLM_SERVE_ARGS"
	;;
*)
	echo "unknown flavor: ${flavor}" >&2
	exit 2
	;;
esac
