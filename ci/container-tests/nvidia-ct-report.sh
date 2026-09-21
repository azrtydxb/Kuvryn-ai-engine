#!/usr/bin/env bash
set -euo pipefail

flavor=${1:?usage: nvidia-ct-report.sh <flavor> <image> [openai-base-url] [model]}
image=${2:?usage: nvidia-ct-report.sh <flavor> <image> [openai-base-url] [model]}
base_url=${3:-}
model=${4:-}

if [[ "${flavor##*-}" != "nvidia" ]]; then
	echo "nvidia CT can only validate *-nvidia flavors, got ${flavor}" >&2
	exit 2
fi

mkdir -p ct-reports
started_at=$(date -u +%Y-%m-%dT%H:%M:%SZ)
ci/container-tests/vendor-device-check.sh nvidia >/tmp/kuvryn-nvidia-device.txt

results='[{"name":"vendor-device","type":"vendor-device","status":"passed","required":true}]'
if [[ -n "$base_url" || -n "$model" ]]; then
	if [[ -z "$base_url" || -z "$model" ]]; then
		echo "openai smoke requires both base URL and model" >&2
		exit 2
	fi
	ci/container-tests/openai-chat-smoke.sh "$base_url" "$model" >/tmp/kuvryn-openai-smoke.txt
	results='[{"name":"vendor-device","type":"vendor-device","status":"passed","required":true},{"name":"openai-chat-smoke","type":"inference-smoke","status":"passed","required":false}]'
fi
ended_at=$(date -u +%Y-%m-%dT%H:%M:%SZ)

cat <<JSON
{"schema":"kuvryn.container-test-report/v1","flavor":"${flavor}","image":"${image}","startedAt":"${started_at}","endedAt":"${ended_at}","results":${results}}
JSON
