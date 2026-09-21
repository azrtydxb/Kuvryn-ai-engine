#!/usr/bin/env bash
set -euo pipefail

flavor=${KUVRYN_DISTRIBUTED_FLAVOR:-${1:-}}
world_size=${WORLD_SIZE:-2}
ranks_ready=${RANKS_READY:-$world_size}
backend=${COLLECTIVE_BACKEND:-unknown}
base_url=${KUVRYN_DISTRIBUTED_BASE_URL:-}
model=${KUVRYN_DISTRIBUTED_MODEL:-}
inference_attested=${KUVRYN_DISTRIBUTED_INFERENCE_PASSED:-}

if [[ -z "$flavor" ]]; then
	echo "usage: KUVRYN_DISTRIBUTED_FLAVOR=<flavor> WORLD_SIZE=<n> RANKS_READY=<n> COLLECTIVE_BACKEND=<backend> distributed-smoke.sh" >&2
	exit 2
fi
if ((world_size < 2)); then
	echo "distributed smoke requires WORLD_SIZE>=2" >&2
	exit 1
fi
if ((ranks_ready != world_size)); then
	echo "distributed smoke requires RANKS_READY to equal WORLD_SIZE" >&2
	exit 1
fi
if [[ "$backend" == unknown || -z "$backend" ]]; then
	echo "COLLECTIVE_BACKEND is required" >&2
	exit 1
fi

inference=false
if [[ -n "$base_url" ]]; then
	if [[ -z "$model" ]]; then
		echo "KUVRYN_DISTRIBUTED_MODEL is required with KUVRYN_DISTRIBUTED_BASE_URL" >&2
		exit 2
	fi
	ci/container-tests/openai-chat-smoke.sh "$base_url" "$model" >/dev/null
	inference=true
else
	inference_attested_normalized=$(printf '%s' "$inference_attested" | tr '[:upper:]' '[:lower:]')
	case "$inference_attested_normalized" in
	1 | true | yes | passed)
		inference=true
		;;
	*)
		echo "distributed smoke requires a real inference request via KUVRYN_DISTRIBUTED_BASE_URL or attested orchestrator evidence via KUVRYN_DISTRIBUTED_INFERENCE_PASSED=true" >&2
		exit 1
		;;
	esac
fi

jq -nc \
	--arg flavor "$flavor" \
	--arg backend "$backend" \
	--argjson worldSize "$world_size" \
	--argjson ranksReady "$ranks_ready" \
	--argjson inference "$inference" \
	'{schema:"kuvryn.distributed-test-report/v1",flavor:$flavor,worldSize:$worldSize,ranksReady:$ranksReady,backend:$backend,inference:$inference}'
