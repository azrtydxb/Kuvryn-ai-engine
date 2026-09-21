package clusterprobe

import "strings"

type Accelerator struct {
	Vendor string `json:"vendor"`
	ID     string `json:"id"`
	Model  string `json:"model"`
}

type FabricInterface struct {
	Name       string `json:"name"`
	Kind       string `json:"kind"`
	RDMADevice string `json:"rdmaDevice,omitempty"`
	PodVisible bool   `json:"podVisible,omitempty"`
}

type HostRecord struct {
	Host         string            `json:"host"`
	ManagementIP string            `json:"managementIp,omitempty"`
	Accelerators []Accelerator     `json:"accelerators"`
	Fabric       []FabricInterface `json:"fabric"`
	Collectives  []string          `json:"collectives"`
	Limitations  []string          `json:"limitations,omitempty"`
}

func DirectProbeScript() string {
	return strings.TrimSpace(`set -eu
printf '===KUVRYN:accelerators===\n'
command -v nvidia-smi >/dev/null 2>&1 && nvidia-smi --query-gpu=index,name --format=csv,noheader || true
command -v rocminfo >/dev/null 2>&1 && rocminfo | grep -E 'Name:.*(gfx|AMD)' || true
command -v xpu-smi >/dev/null 2>&1 && xpu-smi discovery || true
printf '===KUVRYN:network===\n'
ip -o addr show scope global || true
printf '===KUVRYN:fabric===\n'
ls /sys/class/infiniband 2>/dev/null || true
printf '===KUVRYN:runtime===\n'
command -v docker || true
command -v containerd || true
printf '===KUVRYN:end===\n'`) + "\n"
}

func KubernetesDiagnosticCommand(namespace, pod, container string) []string {
	args := []string{"kubectl", "exec", "-n", namespace, pod}
	if container != "" {
		args = append(args, "-c", container)
	}
	args = append(args, "--", "sh", "-ceu", DirectProbeScript())
	return args
}
