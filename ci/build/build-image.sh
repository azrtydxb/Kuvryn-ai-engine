#!/usr/bin/env bash
set -euo pipefail

flavor=${1:?usage: build-image.sh <flavor>}
manifest=${2:-$(ci/lib/manifest-path.sh "$flavor")}
dockerfile=$(awk '/^[[:space:]]+dockerfile:/ { print $2; exit }' "$manifest")
image="${KUVRYN_IMAGE_PREFIX:-ghcr.io/azrtydxb/kuvryn-ai-engine}/${flavor}"
ci_tag=${CI_TAG:-ci}

build_arg_file=$(mktemp)
trap 'rm -f "$build_arg_file"' EXIT
python3 - "$manifest" >"$build_arg_file" <<'PY'
import re
import sys

manifest = sys.argv[1]
in_args = False
with open(manifest, encoding="utf-8") as fh:
    for raw in fh:
        line = raw.rstrip("\n")
        if line == "  args:":
            in_args = True
            continue
        if in_args:
            if not line.startswith("    ") or line.startswith("    -"):
                break
            match = re.match(r"^    ([A-Za-z_][A-Za-z0-9_]*):[ ]*(.*)$", line)
            if not match:
                continue
            key, value = match.groups()
            value = value.strip().strip('"').strip("'")
            print(f"--build-arg\n{key}={value}")
PY
build_args=()
while IFS= read -r flag && IFS= read -r value; do
	build_args+=("$flag" "$value")
done <"$build_arg_file"

cache_dir=".buildx-cache/${flavor}"
cache_args=()
if [[ "${BUILDX_LOCAL_CACHE:-1}" == "1" ]]; then
	driver=$(docker buildx inspect 2>/dev/null | awk '/^Driver:/ { print $2; exit }' || true)
	if [[ "$driver" == "docker" || -z "$driver" ]]; then
		echo "buildx driver ${driver} does not support local cache export; building ${flavor} without buildx local cache" >&2
	else
		mkdir -p "$cache_dir"
		cache_args=(--cache-from "type=local,src=${cache_dir}" --cache-to "type=local,dest=${cache_dir}-next,mode=max")
	fi
fi

docker buildx build \
	--load \
	--label "kuvryn.flavor=${flavor}" \
	--build-arg "BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
	--build-arg "VCS_REF=$(git rev-parse --short HEAD 2>/dev/null || echo unknown)" \
	"${build_args[@]}" \
	"${cache_args[@]}" \
	-t "${image}:ci" \
	-f "$dockerfile" \
	.
if [[ -d "${cache_dir}-next" ]]; then
	rm -rf "$cache_dir"
	mv "${cache_dir}-next" "$cache_dir"
fi

if [[ "${BUILD_PUSH:-0}" == "1" ]]; then
	mkdir -p image-digests
	docker tag "${image}:ci" "${image}:${ci_tag}"
	docker push "${image}:${ci_tag}"
	digest=$(docker image inspect "${image}:${ci_tag}" --format '{{ index .RepoDigests 0 }}' | sed 's/^.*@//')
	if [[ ! "$digest" =~ ^sha256:[0-9a-f]{64}$ ]]; then
		echo "push did not produce an immutable digest for ${image}:${ci_tag}" >&2
		exit 1
	fi
	printf '%s\n' "$digest" >"image-digests/${flavor}.digest"
	printf '%s@%s\n' "$image" "$digest" >"image-digests/${flavor}.ref"
fi
