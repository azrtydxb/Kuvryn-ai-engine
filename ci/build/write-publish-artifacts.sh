#!/usr/bin/env bash
set -euo pipefail

flavor=${1:?usage: write-publish-artifacts.sh <flavor> <digest> <manifest> <report>}
digest=${2:?usage: write-publish-artifacts.sh <flavor> <digest> <manifest> <report>}
manifest=${3:?usage: write-publish-artifacts.sh <flavor> <digest> <manifest> <report>}
report=${4:?usage: write-publish-artifacts.sh <flavor> <digest> <manifest> <report>}

if [[ ! "$digest" =~ ^sha256:[0-9a-f]{64}$ ]]; then
	echo "digest must be an immutable sha256 digest, got ${digest}" >&2
	exit 2
fi

image="${KUVRYN_IMAGE_PREFIX:-ghcr.io/azrtydxb/kuvryn-ai-engine}/${flavor}"
out="publish-artifacts/${flavor}"
mkdir -p "$out"

manifest_sha=$(sha256sum "$manifest" | awk '{print $1}')
report_sha=$(sha256sum "$report" | awk '{print $1}')
revision=$(git rev-parse HEAD 2>/dev/null || echo unknown)
created=$(date -u +%Y-%m-%dT%H:%M:%SZ)
go run ./cmd/kuvryn-publish-tags --manifest "$manifest" --digest "$digest" >"${out}/publish-tags.json"

cat >"${out}/sbom.spdx.json" <<JSON
{
  "spdxVersion": "SPDX-2.3",
  "dataLicense": "CC0-1.0",
  "SPDXID": "SPDXRef-DOCUMENT",
  "name": "${flavor}",
  "documentNamespace": "https://github.com/azrtydxb/kuvryn-ai-engine/sbom/${flavor}/${digest#sha256:}",
  "creationInfo": {
    "created": "${created}",
    "creators": ["Tool: kuvryn-ai-engine-ci"]
  },
  "packages": [
    {
      "name": "${image}",
      "SPDXID": "SPDXRef-Package-${flavor}",
      "downloadLocation": "NOASSERTION",
      "filesAnalyzed": false,
      "checksums": [{"algorithm": "SHA256", "checksumValue": "${digest#sha256:}"}],
      "externalRefs": [{"referenceCategory": "PACKAGE-MANAGER", "referenceType": "purl", "referenceLocator": "pkg:oci/${flavor}@${digest}?repository_url=ghcr.io/azrtydxb/kuvryn-ai-engine/${flavor}"}]
    }
  ]
}
JSON

cat >"${out}/provenance.json" <<JSON
{
  "schema": "kuvryn.provenance/v1",
  "flavor": "${flavor}",
  "image": "${image}",
  "digest": "${digest}",
  "source": "https://github.com/azrtydxb/kuvryn-ai-engine",
  "revision": "${revision}",
  "createdAt": "${created}",
  "publishTags": "${out}/publish-tags.json",
  "materials": [
    {"uri": "${manifest}", "digest": {"sha256": "${manifest_sha}"}},
    {"uri": "${report}", "digest": {"sha256": "${report_sha}"}}
  ]
}
JSON

printf 'wrote publish artifacts to %s\n' "$out"
