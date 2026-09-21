#!/usr/bin/env bash
set -euo pipefail

export DEBIAN_FRONTEND=noninteractive
apt-get update
apt-get install -y --no-install-recommends \
	ca-certificates \
	curl \
	iproute2 \
	jq \
	procps \
	tini
rm -rf /var/lib/apt/lists/*
