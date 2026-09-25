#!/usr/bin/env bash
set -euo pipefail

flavor=${1:?usage: buildctl-build-image.sh <flavor> [image-ref]}
manifest=${KUVRYN_MANIFEST:-$(ci/lib/manifest-path.sh "$flavor")}
image_ref=${2:-${KUVRYN_REGISTRY:-192.168.10.131:5000}/azrtydxb/kuvryn-ai-engine/${flavor}:${KUVRYN_TAG:-kw-test}}
platform=${KUVRYN_PLATFORM:-linux/arm64}
buildkit_host=${BUILDKIT_HOST:-tcp://192.168.10.130:1234}
tls_dir=${BUILDKIT_TLS_DIR:-.buildkit-tls}

if [[ ! -f "$manifest" ]]; then
	echo "manifest not found: $manifest" >&2
	exit 2
fi
if [[ ! -r "$tls_dir/ca.crt" || ! -r "$tls_dir/tls.crt" || ! -r "$tls_dir/tls.key" ]]; then
	echo "missing BuildKit TLS files under $tls_dir; expected ca.crt, tls.crt, tls.key" >&2
	exit 2
fi

dockerfile=$(awk '/^[[:space:]]+dockerfile:/ { print $2; exit }' "$manifest")
if [[ -z "$dockerfile" || ! -f "$dockerfile" ]]; then
	echo "manifest $manifest does not point to a readable Dockerfile" >&2
	exit 2
fi

build_args=()
in_args=0
while IFS= read -r line; do
	if [[ "$line" == "  args:" ]]; then
		in_args=1
		continue
	fi
	if ((in_args)); then
		if [[ ! "$line" =~ ^[[:space:]]{4}[^[:space:]-] ]]; then
			break
		fi
		key=${line%%:*}
		key=${key//[[:space:]]/}
		value=${line#*:}
		value=${value# }
		value=${value%\"}
		value=${value#\"}
		value=${value%\'}
		value=${value#\'}
		if [[ -n "$key" ]]; then
			build_args+=(--opt "build-arg:${key}=${value}")
		fi
	fi
done <"$manifest"

# The runtime wrapper is architecture-specific. Build it for the requested target
# before streaming the context to remote BuildKit.
case "$platform" in
linux/arm64) goarch=arm64 ;;
linux/amd64) goarch=amd64 ;;
*)
	echo "unsupported KUVRYN_PLATFORM for wrapper cross-compile: $platform" >&2
	exit 2
	;;
esac
mkdir -p bin
CGO_ENABLED=0 GOOS=linux GOARCH="$goarch" go build -o bin/kuvryn-runtime-wrapper ./cmd/kuvryn-runtime-wrapper

buildctl \
	--addr "$buildkit_host" \
	--tlscacert "$tls_dir/ca.crt" \
	--tlscert "$tls_dir/tls.crt" \
	--tlskey "$tls_dir/tls.key" \
	build \
	--frontend dockerfile.v0 \
	--local context=. \
	--local dockerfile=. \
	--opt "filename=${dockerfile}" \
	--opt "platform=${platform}" \
	--opt "build-arg:BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
	--opt "build-arg:VCS_REF=$(git rev-parse --short HEAD 2>/dev/null || echo unknown)" \
	"${build_args[@]}" \
	--output "type=image,name=${image_ref},push=true"

mkdir -p image-digests
if command -v crane >/dev/null 2>&1; then
	digest=$(crane digest "$image_ref")
	printf '%s\n' "$digest" >"image-digests/${flavor}.digest"
	printf '%s@%s\n' "${image_ref%:*}" "$digest" >"image-digests/${flavor}.ref"
	printf '%s@%s\n' "${image_ref%:*}" "$digest"
else
	printf '%s\n' "$image_ref"
fi
