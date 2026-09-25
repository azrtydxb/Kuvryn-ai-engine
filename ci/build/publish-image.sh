#!/usr/bin/env bash
set -euo pipefail

flavor=${1:?usage: publish-image.sh <flavor> <digest>}
digest=${2:?usage: publish-image.sh <flavor> <digest>}
manifest=${3:-$(ci/lib/manifest-path.sh "$flavor")}
image="${KUVRYN_IMAGE_PREFIX:-ghcr.io/azrtydxb/kuvryn-ai-engine}/${flavor}"

if [[ ! "$digest" =~ ^sha256:[0-9a-f]{64}$ ]]; then
	echo "digest must be an immutable sha256 digest, got ${digest}" >&2
	exit 2
fi

report="ct-reports/${flavor}.json"
if [[ ! -f "$report" ]]; then
	echo "missing required CT report for ${flavor}" >&2
	exit 1
fi

go run ./cmd/kuvryn-manifest-check --root .
go run ./cmd/kuvryn-ct-validate --manifest "$manifest" --report "$report"
ci/build/write-publish-artifacts.sh "$flavor" "$digest" "$manifest" "$report"

mapfile -t tag_args < <(go run ./cmd/kuvryn-publish-tags --manifest "$manifest" --digest "$digest" --format args)
docker buildx imagetools create \
	"${tag_args[@]}" \
	"${image}@${digest}"

mkdir -p lock-updates
go run ./cmd/kuvryn-lock-accept \
	--manifest "$manifest" \
	--report "$report" \
	--digest "$digest" \
	--output "lock-updates/${flavor}.upstream.lock"
