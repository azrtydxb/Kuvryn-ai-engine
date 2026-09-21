#!/usr/bin/env bash
set -euo pipefail

version=${VLLM_VERSION:?VLLM_VERSION is required}
ray_version=${RAY_VERSION:?RAY_VERSION is required}
profile=${KUVRYN_VENDOR:?KUVRYN_VENDOR is required}

case "$profile" in
nvidia)
	python3 -m pip install --no-cache-dir "ray[default]==${ray_version}" "vllm==${version}"
	;;
amd)
	# ROCm images must stay CUDA-free; use the ROCm PyTorch index before installing vLLM.
	rocm_index=${PYTORCH_ROCM_INDEX:-https://download.pytorch.org/whl/rocm6.2}
	python3 -m pip install --no-cache-dir --index-url "$rocm_index" torch torchvision torchaudio
	python3 -m pip install --no-cache-dir "ray[default]==${ray_version}" "vllm==${version}"
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
