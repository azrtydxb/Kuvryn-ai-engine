package placement

import (
	"fmt"
	"time"

	"github.com/azrtydxb/kuvryn-ai-engine/internal/clusterprobe"
	"github.com/azrtydxb/kuvryn-ai-engine/internal/engineimage"
)

type DeploymentRequest struct {
	Engine              string
	Vendor              string
	Distributed         bool
	RDMARequired        bool
	Kubernetes          bool
	RDMAValidation      *clusterprobe.RDMAValidation
	ExplicitMixedVendor bool
}

type PlacementPlan struct {
	Flavor string
	Hosts  []string
	Env    map[string]string
}

func CheckPlacement(req DeploymentRequest, inventory []clusterprobe.HostRecord, flavor engineimage.ImageManifest) (PlacementPlan, error) {
	if req.Engine != flavor.Engine || req.Vendor != flavor.Vendor {
		return PlacementPlan{}, fmt.Errorf("requested %s/%s incompatible with flavor %s/%s", req.Engine, req.Vendor, flavor.Engine, flavor.Vendor)
	}
	if req.Distributed && !flavor.Capabilities.Distributed {
		return PlacementPlan{}, fmt.Errorf("flavor %s is not distributed-capable", flavor.Name)
	}
	vendors := map[string]bool{}
	for _, h := range inventory {
		for _, a := range h.Accelerators {
			vendors[a.Vendor] = true
		}
	}
	if len(vendors) > 1 && !req.ExplicitMixedVendor {
		return PlacementPlan{}, fmt.Errorf("mixed-vendor auto-placement requires explicit layout")
	}
	if req.RDMARequired {
		if req.RDMAValidation == nil {
			return PlacementPlan{}, fmt.Errorf("RDMA-required placement needs successful RDMA validation")
		}
		validation := *req.RDMAValidation
		validation.Required = true
		if !validation.UsableForRequired(time.Now(), time.Hour, req.Kubernetes) {
			return PlacementPlan{}, fmt.Errorf("RDMA validation is not usable for required placement: %s", validation.Message)
		}
	}
	hosts := make([]string, 0, len(inventory))
	for _, h := range inventory {
		hosts = append(hosts, h.Host)
	}
	env, err := PlanCollectiveEnv(inventory, flavor)
	if err != nil {
		return PlacementPlan{}, err
	}
	return PlacementPlan{Flavor: flavor.Name, Hosts: hosts, Env: env.Env}, nil
}
