#!/usr/bin/env bash
set -euo pipefail

flavor=${1:?usage: pull-smoke.sh <flavor> <image>}
image=${2:?usage: pull-smoke.sh <flavor> <image>}
namespace=${KUVRYN_K8S_NAMESPACE:-default}
context=${KUVRYN_KUBECTL_CONTEXT:-}

kube() {
	if [[ -n "$context" ]]; then
		kubectl --context "$context" "$@"
	else
		kubectl "$@"
	fi
}

nodes=$(kube get nodes -o jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}')
for node in $nodes; do
	pod="kuvryn-pull-${flavor}-${node//[^a-zA-Z0-9-]/-}"
	pod=${pod:0:63}
	kube -n "$namespace" delete pod "$pod" --ignore-not-found --wait=true >/dev/null 2>&1 || true
	kube -n "$namespace" run "$pod" \
		--image="$image" \
		--restart=Never \
		--overrides="{\"spec\":{\"nodeName\":\"$node\",\"tolerations\":[{\"operator\":\"Exists\"}],\"containers\":[{\"name\":\"$pod\",\"image\":\"$image\",\"imagePullPolicy\":\"Always\",\"command\":[\"sh\",\"-c\",\"echo pulled ${flavor} on ${node}\"]}]}}" >/dev/null
	kube -n "$namespace" wait --for=jsonpath='{.status.phase}'=Succeeded "pod/$pod" --timeout=180s >/dev/null
	kube -n "$namespace" logs "$pod"
	kube -n "$namespace" delete pod "$pod" --ignore-not-found >/dev/null
done
