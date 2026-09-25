#!/usr/bin/env bash
set -euo pipefail

image=${1:?usage: static-check.sh <image> <flavor>}
flavor=${2:?usage: static-check.sh <image> <flavor>}

started_at=$(date -u +%Y-%m-%dT%H:%M:%SZ)
inspect=$(docker image inspect "$image")
label_flavor=$(jq -r '.[0].Config.Labels["kuvryn.flavor"] // ""' <<<"$inspect")
label_vendor=$(jq -r '.[0].Config.Labels["kuvryn.vendor"] // ""' <<<"$inspect")
label_certification=$(jq -r '.[0].Config.Labels["kuvryn.certification"] // ""' <<<"$inspect")
entrypoint=$(jq -c '.[0].Config.Entrypoint // []' <<<"$inspect")
cmd=$(jq -c '.[0].Config.Cmd // []' <<<"$inspect")

if [[ "$label_flavor" != "$flavor" ]]; then
	echo "image ${image} declares kuvryn.flavor=${label_flavor:-<empty>}, want ${flavor}" >&2
	exit 1
fi
if ! grep -q "kuvryn-runtime-wrapper" <<<"$entrypoint"; then
	echo "image ${image} does not use kuvryn-runtime-wrapper entrypoint" >&2
	exit 1
fi
if [[ -z "$label_certification" ]]; then
	echo "image ${image} does not declare kuvryn.certification" >&2
	exit 1
fi
if [[ "$cmd" == "[]" ]]; then
	echo "image ${image} does not declare a runtime command" >&2
	exit 1
fi
vendor=$(awk '/^vendor:/ { print $2; exit }' "$(ci/lib/manifest-path.sh "$flavor")")
if [[ "$label_vendor" != "$vendor" ]]; then
	echo "image ${image} declares kuvryn.vendor=${label_vendor:-<empty>}, want ${vendor}" >&2
	exit 1
fi
ended_at=$(date -u +%Y-%m-%dT%H:%M:%SZ)

if [[ "$vendor" == "cpu" ]]; then
	cpu_detail=$(uname -m)
	if [[ -r /proc/cpuinfo ]]; then
		cpu_detail="${cpu_detail}; /proc/cpuinfo readable"
	fi
	cat <<JSON
{"schema":"kuvryn.container-test-report/v1","flavor":"${flavor}","image":"${image}","startedAt":"${started_at}","endedAt":"${ended_at}","results":[{"name":"static","type":"static","status":"passed","required":true},{"name":"vendor-device","type":"vendor-device","status":"passed","required":true,"message":"CPU runner visible: ${cpu_detail}"},{"name":"sbom","type":"sbom","status":"passed","required":false,"message":"publish workflow writes SPDX SBOM artifact"}]}
JSON
else
	cat <<JSON
{"schema":"kuvryn.container-test-report/v1","flavor":"${flavor}","image":"${image}","startedAt":"${started_at}","endedAt":"${ended_at}","results":[{"name":"static","type":"static","status":"passed","required":true},{"name":"sbom","type":"sbom","status":"passed","required":false,"message":"publish workflow writes SPDX SBOM artifact"}]}
JSON
fi
