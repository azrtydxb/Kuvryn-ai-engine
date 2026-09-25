#!/usr/bin/env bash
set -euo pipefail

version=${VLLM_VERSION:?VLLM_VERSION is required}
ray_version=${RAY_VERSION:?RAY_VERSION is required}
profile=${KUVRYN_VENDOR:?KUVRYN_VENDOR is required}

install_rocm_vllm() {
	# PyPI only ships CUDA builds of vLLM, and vLLM publishes ROCm wheels only for the ROCm
	# it pins (7.2.x). To run on the current ROCm, build vLLM from source against AMD's
	# ROCm 10 wheels; the rocm[devel] extra provides hipcc and headers, torch pulls the
	# matching rocm[libraries].
	local rocm_version torch_version torchvision_version triton_version aiter_ref rocm_index site
	rocm_version=${ROCM_VERSION:?ROCM_VERSION is required}
	torch_version=${TORCH_VERSION:?TORCH_VERSION is required}
	torchvision_version=${TORCHVISION_VERSION:?TORCHVISION_VERSION is required}
	triton_version=${TRITON_VERSION:?TRITON_VERSION is required}
	aiter_ref=${AITER_REF:?AITER_REF is required}
	rocm_index=${ROCM_WHEEL_INDEX:-https://stable.repo.amd.com/rocm/whl-next/}

	export DEBIAN_FRONTEND=noninteractive
	apt-get update
	apt-get install -y --no-install-recommends \
		g++ \
		git \
		libdrm-dev \
		libnuma-dev \
		pkg-config \
		protobuf-compiler
	rm -rf /var/lib/apt/lists/*

	python3 -m pip install --no-cache-dir --index-url "$rocm_index" --extra-index-url https://pypi.org/simple \
		"rocm[libraries,devel,device-all]==${rocm_version}" \
		"torch[device-all]==${torch_version}" \
		"torchvision==${torchvision_version}" \
		"triton==${triton_version}"
	rocm-sdk init
	site=$(python3 -c 'import sysconfig; print(sysconfig.get_paths()["purelib"])')
	ln -sfn "${site}/_rocm_sdk_devel" /opt/rocm
	python3 -m pip install --no-cache-dir packaging 'cmake<4' ninja wheel 'setuptools<80' setuptools-scm setuptools-rust pybind11 Cython
	python3 -m pip install --no-cache-dir "${site}/_rocm_sdk_core/share/amd_smi"

	# vLLM's Rust frontend needs a toolchain only while building.
	export CARGO_HOME=/tmp/cargo RUSTUP_HOME=/tmp/rustup
	export PATH="/tmp/cargo/bin:${PATH}"
	curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y --default-toolchain none --no-modify-path

	# AITER kernels are JIT-compiled on first use; on RDNA the Triton kernels run and the
	# CDNA-only CK/ASM kernels are skipped, so nothing is prebuilt here.
	git clone --recursive --depth 1 --branch "$aiter_ref" https://github.com/ROCm/aiter.git /tmp/aiter
	(cd /tmp/aiter &&
		python3 -m pip install --no-cache-dir -r requirements.txt &&
		AITER_USE_SYSTEM_TRITON=1 PREBUILD_KERNELS=0 python3 -m pip install --no-cache-dir --no-build-isolation .)

	git clone --depth 1 --branch "v${version}" https://github.com/vllm-project/vllm.git /tmp/vllm
	(cd /tmp/vllm &&
		rustup toolchain install "$(sed -n 's/^channel *= *"\(.*\)"/\1/p' rust-toolchain.toml)" &&
		python3 -m pip install --no-cache-dir -r requirements/rocm.txt &&
		VLLM_TARGET_DEVICE=rocm python3 -m pip install --no-cache-dir --no-build-isolation .)
	python3 -m pip install --no-cache-dir "ray[default]==${ray_version}"
	rm -rf /tmp/aiter /tmp/vllm /tmp/cargo /tmp/rustup
}

case "$profile" in
nvidia)
	python3 -m pip install --no-cache-dir "ray[default]==${ray_version}" "vllm==${version}"
	;;
amd)
	install_rocm_vllm
	;;
intel)
	# Intel support is build-only until hardware CT exists; keep dependencies in the Intel image only.
	python3 -m pip install --no-cache-dir "ray[default]==${ray_version}" "vllm==${version}"
	;;
*)
	echo "unsupported vLLM vendor profile: ${profile}" >&2
	exit 2
	;;
esac

command -v vllm >/dev/null
python3 - <<'PY'
import importlib
importlib.import_module("vllm")
PY
