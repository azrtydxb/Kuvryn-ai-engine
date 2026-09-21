#!/usr/bin/env bash
set -euo pipefail

flavor=${1:?usage: nvidia-image-check.sh <flavor> <image>}
image=${2:?usage: nvidia-image-check.sh <flavor> <image>}

if [[ "${flavor##*-}" != "nvidia" ]]; then
	echo "nvidia image check can only validate *-nvidia flavors, got ${flavor}" >&2
	exit 2
fi

started_at=$(date -u +%Y-%m-%dT%H:%M:%SZ)
static_json=$(ci/container-tests/static-check.sh "$image" "$flavor")

vendor_status=passed
vendor_message="NVIDIA devices visible inside container"
if ! docker run --rm --gpus all --entrypoint /bin/sh "$image" -lc 'test -e /dev/nvidiactl || test -e /dev/nvidia0 || test -d /proc/driver/nvidia'; then
	vendor_status=failed
	vendor_message="NVIDIA devices were not visible inside container"
fi

wrapper_status=passed
wrapper_message="kuvryn-runtime-wrapper executed /bin/true"
if ! docker run --rm --gpus all "$image" /bin/true; then
	wrapper_status=failed
	wrapper_message="kuvryn-runtime-wrapper failed to execute /bin/true"
fi

runtime_status=passed
runtime_message="runtime command exists"
case "$flavor" in
vllm-nvidia)
	runtime_probe='command -v vllm || python -m vllm --help >/dev/null'
	;;
sglang-nvidia)
	runtime_probe='command -v sglang || python -m sglang.launch_server --help >/dev/null'
	;;
llama-cpp-nvidia)
	runtime_probe='command -v llama-server'
	;;
tensorrt-llm-nvidia)
	runtime_probe='command -v trtllm-serve'
	;;
*)
	runtime_probe='true'
	;;
esac
if ! docker run --rm --gpus all --entrypoint /bin/sh "$image" -lc "$runtime_probe" >/tmp/kuvryn-runtime-command.txt 2>&1; then
	runtime_status=failed
	runtime_message="runtime command for ${flavor} is missing or not runnable"
fi

ended_at=$(date -u +%Y-%m-%dT%H:%M:%SZ)
jq -c \
	--arg started "$started_at" \
	--arg ended "$ended_at" \
	--arg vendor_status "$vendor_status" \
	--arg vendor_message "$vendor_message" \
	--arg wrapper_status "$wrapper_status" \
	--arg wrapper_message "$wrapper_message" \
	--arg runtime_status "$runtime_status" \
	--arg runtime_message "$runtime_message" \
	'.startedAt=$started | .endedAt=$ended | .results += [
		{"name":"vendor-device","type":"vendor-device","status":$vendor_status,"required":true,"message":$vendor_message},
		{"name":"wrapper-exec","type":"runtime-wrapper","status":$wrapper_status,"required":true,"message":$wrapper_message},
		{"name":"runtime-command","type":"runtime-command","status":$runtime_status,"required":true,"message":$runtime_message}
	]' <<<"$static_json"
