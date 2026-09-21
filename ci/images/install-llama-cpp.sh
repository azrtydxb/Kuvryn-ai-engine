#!/usr/bin/env bash
set -euo pipefail

ref=${LLAMA_CPP_REF:?LLAMA_CPP_REF is required}
repo=${LLAMA_CPP_REPO:-https://github.com/ggerganov/llama.cpp.git}
build_dir=${LLAMA_CPP_BUILD_DIR:-/tmp/llama.cpp}
cmake_flags=${LLAMA_CPP_CMAKE_FLAGS:-}

export DEBIAN_FRONTEND=noninteractive
apt-get update
apt-get install -y --no-install-recommends \
	build-essential \
	cmake \
	git \
	libcurl4-openssl-dev \
	pkg-config
rm -rf /var/lib/apt/lists/*

rm -rf "$build_dir"
git clone --depth 1 --branch "$ref" "$repo" "$build_dir"
# shellcheck disable=SC2086 # cmake flags are intentionally split by the Dockerfile per flavor.
cmake -S "$build_dir" -B "$build_dir/build" \
	-DCMAKE_BUILD_TYPE=Release \
	-DLLAMA_BUILD_SERVER=ON \
	-DLLAMA_CURL=ON \
	$cmake_flags
cmake --build "$build_dir/build" --target llama-server -j"$(nproc)"
find "$build_dir/build/bin" -maxdepth 1 -type f -name '*.so*' -exec install -m 0755 {} /usr/local/lib/ \;
install -m 0755 "$build_dir/build/bin/llama-server" /usr/local/bin/llama-server
ldconfig
if [[ "${LLAMA_CPP_SKIP_HELP:-0}" != "1" ]]; then
	llama-server --help >/dev/null
fi
rm -rf "$build_dir"
