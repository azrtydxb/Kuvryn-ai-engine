#!/usr/bin/env bash
set -euo pipefail

registry=${1:-192.168.10.131:5000}
ca_file=${KUVRYN_K3S_REGISTRY_CA_FILE:-/etc/rancher/k3s/cluster-ca.crt}
context=${KUVRYN_KUBECTL_CONTEXT:-}
namespace=${KUVRYN_K8S_NAMESPACE:-default}

kube() {
	if [[ -n "$context" ]]; then
		kubectl --context "$context" "$@"
	else
		kubectl "$@"
	fi
}

if [[ "$registry" != *:* ]]; then
	echo "registry must include host:port, got $registry" >&2
	exit 2
fi

host_dir="/var/lib/rancher/k3s/agent/etc/containerd/certs.d/${registry}"
hosts_toml=$(
	cat <<EOF
server = "https://${registry}/v2"
capabilities = ["pull", "resolve", "push"]
ca = ["${ca_file}"]

[host]
EOF
)

nodes=$(kube get nodes -o jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}')
for node in $nodes; do
	pod="kw-registry-ca-${node//[^a-zA-Z0-9-]/-}"
	kube -n "$namespace" delete pod "$pod" --ignore-not-found --wait=false >/dev/null 2>&1 || true
	kube -n "$namespace" apply -f - >/dev/null <<YAML
apiVersion: v1
kind: Pod
metadata:
  name: ${pod}
spec:
  restartPolicy: Never
  nodeName: ${node}
  tolerations:
    - operator: Exists
  containers:
    - name: fix
      image: busybox:1.36
      command:
        - sh
        - -c
        - |
          set -eu
          dir=/host${host_dir}
          mkdir -p "\$dir"
          cat > "\$dir/hosts.toml" <<'EOF'
${hosts_toml}
EOF
          chmod 0644 "\$dir/hosts.toml"
          echo "${node}: wrote ${host_dir}/hosts.toml"
      securityContext:
        privileged: true
      volumeMounts:
        - name: host
          mountPath: /host
  volumes:
    - name: host
      hostPath:
        path: /
YAML
	kube -n "$namespace" wait --for=jsonpath='{.status.phase}'=Succeeded "pod/$pod" --timeout=90s >/dev/null
	kube -n "$namespace" logs "$pod"
	kube -n "$namespace" delete pod "$pod" --ignore-not-found >/dev/null
done
