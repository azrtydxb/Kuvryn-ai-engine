#!/usr/bin/env bash
set -euo pipefail

version=${SGLANG_VERSION:?SGLANG_VERSION is required}
profile=${KUVRYN_VENDOR:?KUVRYN_VENDOR is required}

case "$profile" in
nvidia)
	env PIP_CONSTRAINT= python3 -m pip install --no-cache-dir "sglang[all]==${version}"
	;;
amd)
	rocm_index=${PYTORCH_ROCM_INDEX:-https://download.pytorch.org/whl/rocm6.2}
	python3 -m pip install --no-cache-dir --index-url "$rocm_index" torch torchvision torchaudio
	python3 -m pip install --no-cache-dir "sglang==${version}"
	;;
intel)
	python3 -m pip install --no-cache-dir "sglang==${version}"
	;;
*)
	echo "unsupported SGLang vendor profile: ${profile}" >&2
	exit 2
	;;
esac

python3 - <<'PY'
import importlib
importlib.import_module("sglang")
PY
cat >/usr/local/bin/sglang <<'SH'
#!/usr/bin/env sh
exec python3 -m sglang.launch_server "$@"
SH
chmod 0755 /usr/local/bin/sglang
command -v sglang >/dev/null
