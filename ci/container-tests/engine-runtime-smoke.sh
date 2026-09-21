#!/usr/bin/env bash
set -euo pipefail

flavor=${1:?usage: engine-runtime-smoke.sh <flavor> <image>}
image=${2:?usage: engine-runtime-smoke.sh <flavor> <image>}
model=${KUVRYN_SMOKE_MODEL:-}
port=${KUVRYN_SMOKE_PORT:-8000}
env_file=${KUVRYN_SMOKE_ENV_FILE:-}
model_mount=${KUVRYN_SMOKE_MODEL_MOUNT:-}
container="kuvryn-smoke-${flavor}-$$"

if [[ -n "$env_file" ]]; then
	# shellcheck disable=SC1090 # env file is caller-provided fixture output.
	source "$env_file"
	model=${KUVRYN_SMOKE_MODEL:-$model}
	model_mount=${KUVRYN_SMOKE_MODEL_MOUNT:-$model_mount}
fi
if [[ -z "$model" && "$flavor" != tensorrt-llm-* ]]; then
	echo "KUVRYN_SMOKE_MODEL must point to a real tiny model fixture for ${flavor}; run prepare-smoke-fixture.sh first" >&2
	exit 2
fi

runtime_args=()
case "$flavor" in
vllm-*)
	# shellcheck disable=SC2206 # caller intentionally supplies vLLM smoke-only serve flags.
	vllm_extra_args=(${KUVRYN_VLLM_SERVE_ARGS:-})
	runtime_args=(vllm serve "$model" --host 0.0.0.0 --port "$port" "${vllm_extra_args[@]}")
	;;
sglang-*)
	# shellcheck disable=SC2206 # caller intentionally supplies SGLang smoke-only serve flags.
	sglang_extra_args=(${KUVRYN_SGLANG_SERVE_ARGS:-})
	runtime_args=(sglang serve --model-path "$model" --host 0.0.0.0 --port "$port" "${sglang_extra_args[@]}")
	;;
llama-cpp-*)
	# shellcheck disable=SC2206 # caller intentionally supplies llama.cpp smoke-only serve flags.
	llama_cpp_extra_args=(${KUVRYN_LLAMA_CPP_SERVE_ARGS:-})
	runtime_args=(llama-server -m "$model" --host 0.0.0.0 --port "$port" "${llama_cpp_extra_args[@]}")
	;;
tensorrt-llm-*)
	# TensorRT-LLM serves an already-built engine/checkpoint. Pass the exact tested args,
	# for example: KUVRYN_TRTLLM_SERVE_ARGS="/models/tiny-trtllm --host 0.0.0.0 --port 8000".
	if [[ -z "${KUVRYN_TRTLLM_SERVE_ARGS:-}" ]]; then
		echo "KUVRYN_TRTLLM_SERVE_ARGS is required for TensorRT-LLM smoke" >&2
		exit 2
	fi
	# shellcheck disable=SC2206 # caller intentionally supplies trtllm-serve arguments.
	runtime_args=(trtllm-serve ${KUVRYN_TRTLLM_SERVE_ARGS})
	;;
*)
	echo "unknown flavor: ${flavor}" >&2
	exit 2
	;;
esac

cleanup() {
	docker rm -f "$container" >/dev/null 2>&1 || true
}
trap cleanup EXIT

docker_args=(-d --rm --name "$container")
if [[ -n "$env_file" ]]; then
	docker_args+=(--env-file "$env_file")
fi
if [[ -n "$model_mount" ]]; then
	docker_args+=(-v "$model_mount")
fi
# shellcheck disable=SC2206 # KUVRYN_DOCKER_RUN_FLAGS intentionally allows caller-supplied docker flags such as --gpus all.
extra_docker_args=(${KUVRYN_DOCKER_RUN_FLAGS:-})
docker_args+=("${extra_docker_args[@]}" -p "127.0.0.1:${port}:${port}")
docker run "${docker_args[@]}" "$image" "${runtime_args[@]}" >/dev/null

base_url="http://127.0.0.1:${port}"
deadline=$((SECONDS + ${KUVRYN_SMOKE_TIMEOUT_SECONDS:-180}))
until curl -fsS "${base_url}/v1/models" >/dev/null 2>&1 || curl -fsS "${base_url}/health" >/dev/null 2>&1; do
	if ((SECONDS >= deadline)); then
		docker logs "$container" >&2 || true
		echo "runtime did not become ready before timeout" >&2
		exit 1
	fi
	sleep 2
done

if [[ "$flavor" == llama-cpp-* ]]; then
	ci/container-tests/openai-chat-smoke.sh "$base_url" "${KUVRYN_SMOKE_MODEL_NAME:-$model}"
else
	ci/container-tests/openai-chat-smoke.sh "$base_url" "${KUVRYN_SMOKE_MODEL_NAME:-$model}"
fi
