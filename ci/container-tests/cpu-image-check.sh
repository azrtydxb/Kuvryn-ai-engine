#!/usr/bin/env bash
set -euo pipefail

flavor=${1:?usage: cpu-image-check.sh <flavor> <image>}
image=${2:?usage: cpu-image-check.sh <flavor> <image>}

if [[ "${flavor##*-}" != "cpu" ]]; then
	echo "cpu image check can only validate *-cpu flavors, got ${flavor}" >&2
	exit 2
fi

started_at=$(date -u +%Y-%m-%dT%H:%M:%SZ)
static_json=$(ci/container-tests/static-check.sh "$image" "$flavor")

wrapper_status=passed
wrapper_message="kuvryn-runtime-wrapper executed /bin/true"
if ! docker run --rm "$image" /bin/true; then
	wrapper_status=failed
	wrapper_message="kuvryn-runtime-wrapper failed to execute /bin/true"
fi

runtime_status=passed
runtime_message="runtime command exists"
case "$flavor" in
llama-cpp-cpu)
	runtime_probe='command -v llama-server && llama-server --help >/dev/null'
	;;
*)
	runtime_probe='true'
	;;
esac
if ! docker run --rm --entrypoint /bin/sh "$image" -lc "$runtime_probe" >/tmp/kuvryn-runtime-command.txt 2>&1; then
	runtime_status=failed
	runtime_message="runtime command for ${flavor} is missing or not runnable"
fi

ended_at=$(date -u +%Y-%m-%dT%H:%M:%SZ)
jq -c \
	--arg started "$started_at" \
	--arg ended "$ended_at" \
	--arg wrapper_status "$wrapper_status" \
	--arg wrapper_message "$wrapper_message" \
	--arg runtime_status "$runtime_status" \
	--arg runtime_message "$runtime_message" \
	'.startedAt=$started | .endedAt=$ended | .results += [
		{"name":"wrapper-exec","type":"runtime-wrapper","status":$wrapper_status,"required":true,"message":$wrapper_message},
		{"name":"runtime-command","type":"runtime-command","status":$runtime_status,"required":true,"message":$runtime_message}
	]' <<<"$static_json"
