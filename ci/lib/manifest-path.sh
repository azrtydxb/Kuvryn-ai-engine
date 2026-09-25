#!/usr/bin/env bash
# Print the image.yaml path for a flavor by its declared name.
# Flavors may carry a variant suffix (engine-vendor-variant), so the path is
# looked up rather than derived by splitting the name.
set -euo pipefail

flavor=${1:?usage: manifest-path.sh <flavor>}
for manifest in engines/*/*/image.yaml; do
	if grep -qx "name: ${flavor}" "$manifest"; then
		printf '%s\n' "$manifest"
		exit 0
	fi
done
echo "no manifest declares name: ${flavor}" >&2
exit 2
