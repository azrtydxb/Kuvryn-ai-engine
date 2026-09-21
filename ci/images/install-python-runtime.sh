#!/usr/bin/env bash
set -euo pipefail

export DEBIAN_FRONTEND=noninteractive
apt-get update
apt-get install -y --no-install-recommends \
	python3 \
	python3-pip \
	python3-venv
rm -rf /var/lib/apt/lists/*
python3 -m venv "${VIRTUAL_ENV:-/opt/kuvryn/venv}"
"${VIRTUAL_ENV:-/opt/kuvryn/venv}/bin/python" -m pip install --no-cache-dir --upgrade pip setuptools wheel
