package placement

import (
	"testing"
	"time"

	"github.com/azrtydxb/kuvryn-ai-engine/internal/clusterprobe"
	"github.com/azrtydxb/kuvryn-ai-engine/internal/engineimage"
)

func TestPlanCollectiveEnvNCCL(t *testing.T) {
	flavor := engineimage.ImageManifest{Name: "vllm-nvidia", Engine: "vllm", Vendor: "nvidia", Capabilities: engineimage.Capabilities{Collective: "nccl", Distributed: true}}
	plan, err := PlanCollectiveEnv([]clusterprobe.HostRecord{{Fabric: []clusterprobe.FabricInterface{{Name: "mlx5_0"}}}}, flavor)
	if err != nil {
		t.Fatal(err)
	}
	if plan.Env["NCCL_IB_HCA"] != "mlx5_0" {
		t.Fatalf("NCCL_IB_HCA = %q", plan.Env["NCCL_IB_HCA"])
	}
}

func TestPlacementRefusesMixedVendor(t *testing.T) {
	flavor := engineimage.ImageManifest{Name: "vllm-nvidia", Engine: "vllm", Vendor: "nvidia", Capabilities: engineimage.Capabilities{Collective: "nccl", Distributed: true}}
	inventory := []clusterprobe.HostRecord{{Host: "a", Accelerators: []clusterprobe.Accelerator{{Vendor: "nvidia"}}}, {Host: "b", Accelerators: []clusterprobe.Accelerator{{Vendor: "amd"}}}}
	_, err := CheckPlacement(DeploymentRequest{Engine: "vllm", Vendor: "nvidia", Distributed: true}, inventory, flavor)
	if err == nil {
		t.Fatal("expected mixed vendor error")
	}
}

func TestPlacementRequiresFreshRDMAValidation(t *testing.T) {
	flavor := engineimage.ImageManifest{Name: "vllm-nvidia", Engine: "vllm", Vendor: "nvidia", Capabilities: engineimage.Capabilities{Collective: "nccl", Distributed: true}}
	inventory := []clusterprobe.HostRecord{{Host: "a", Accelerators: []clusterprobe.Accelerator{{Vendor: "nvidia"}}, Fabric: []clusterprobe.FabricInterface{{Name: "mlx5_0", PodVisible: true}}}}
	_, err := CheckPlacement(DeploymentRequest{Engine: "vllm", Vendor: "nvidia", Distributed: true, RDMARequired: true}, inventory, flavor)
	if err == nil {
		t.Fatal("expected missing RDMA validation error")
	}
	validation := &clusterprobe.RDMAValidation{Reachable: true, TransferOK: true, CollectiveOK: true, PodVisible: true, CheckedAt: time.Now()}
	if _, err := CheckPlacement(DeploymentRequest{Engine: "vllm", Vendor: "nvidia", Distributed: true, RDMARequired: true, RDMAValidation: validation}, inventory, flavor); err != nil {
		t.Fatalf("validated RDMA placement failed: %v", err)
	}
	validation.PodVisible = false
	_, err = CheckPlacement(DeploymentRequest{Engine: "vllm", Vendor: "nvidia", Distributed: true, RDMARequired: true, Kubernetes: true, RDMAValidation: validation}, inventory, flavor)
	if err == nil {
		t.Fatal("expected pod-invisible RDMA validation error")
	}
}
