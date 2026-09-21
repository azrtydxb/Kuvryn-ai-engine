package clusterprobe

import (
	"testing"
	"time"
)

func TestParseSections(t *testing.T) {
	sections, err := ParseSections("===KUVRYN:accelerators===\n0,A100\n===KUVRYN:network===\neth0\n===KUVRYN:fabric===\nmlx5_0\n===KUVRYN:runtime===\ndocker\n===KUVRYN:end===\n")
	if err != nil {
		t.Fatal(err)
	}
	if sections["fabric"][0] != "mlx5_0" {
		t.Fatalf("fabric = %#v", sections["fabric"])
	}
	if _, err := ParseSections("===KUVRYN:accelerators===\n"); err == nil {
		t.Fatal("expected malformed probe error")
	}
}

func TestKubernetesDiagnosticCommand(t *testing.T) {
	cmd := KubernetesDiagnosticCommand("default", "probe-pod", "engine")
	if cmd[0] != "kubectl" || cmd[2] != "-n" || cmd[3] != "default" || cmd[4] != "probe-pod" {
		t.Fatalf("command = %#v", cmd)
	}
	if cmd[len(cmd)-2] != "-ceu" {
		t.Fatalf("command missing shell probe: %#v", cmd)
	}
}

func TestTopologyAndRDMA(t *testing.T) {
	records := []HostRecord{{Host: "a", Fabric: []FabricInterface{{Name: "mlx5_0", PodVisible: true}}}, {Host: "b", Fabric: []FabricInterface{{Name: "mlx5_1", PodVisible: true}}}}
	if got := ClassifyTopology(records, true); got != TopologyKubernetesFabric {
		t.Fatalf("topology = %s", got)
	}
	v := RDMAValidation{Required: true, Reachable: true, TransferOK: true, CollectiveOK: true, PodVisible: true, CheckedAt: time.Now()}
	if !v.UsableForRequired(time.Now(), time.Minute, true) {
		t.Fatal("RDMA should be usable")
	}
	v.PodVisible = false
	if v.UsableForRequired(time.Now(), time.Minute, true) {
		t.Fatal("pod-invisible RDMA should fail")
	}
}
