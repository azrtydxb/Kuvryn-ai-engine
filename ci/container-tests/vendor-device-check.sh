#!/usr/bin/env bash
set -euo pipefail

vendor=${1:?usage: vendor-device-check.sh <cpu|nvidia|amd|intel>}
case "$vendor" in
cpu) test -r /proc/cpuinfo ;;
nvidia) command -v nvidia-smi >/dev/null && nvidia-smi -L ;;
amd) test -e /dev/kfd || command -v rocminfo >/dev/null ;;
intel) test -e /dev/dri || command -v xpu-smi >/dev/null ;;
*)
	echo "unsupported vendor ${vendor}" >&2
	exit 2
	;;
esac
