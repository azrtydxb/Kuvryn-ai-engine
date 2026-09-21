package placement

import (
	"fmt"
	"strings"

	"github.com/azrtydxb/kuvryn-ai-engine/internal/clusterprobe"
	"github.com/azrtydxb/kuvryn-ai-engine/internal/engineimage"
)

type CollectiveEnvPlan struct {
	Backend string
	Env     map[string]string
}

func PlanCollectiveEnv(records []clusterprobe.HostRecord, flavor engineimage.ImageManifest) (CollectiveEnvPlan, error) {
	env := map[string]string{}
	ifs := fabricNames(records)
	switch flavor.Capabilities.Collective {
	case "none", "":
		return CollectiveEnvPlan{Backend: "none", Env: env}, nil
	case "nccl":
		if flavor.Vendor != "nvidia" {
			return CollectiveEnvPlan{}, fmt.Errorf("NCCL requires nvidia flavor")
		}
		env["NCCL_SOCKET_IFNAME"] = strings.Join(ifs, ",")
		if len(ifs) > 0 {
			env["NCCL_IB_HCA"] = strings.Join(ifs, ",")
		}
	case "rccl":
		if flavor.Vendor != "amd" {
			return CollectiveEnvPlan{}, fmt.Errorf("RCCL requires amd flavor")
		}
		env["RCCL_SOCKET_IFNAME"] = strings.Join(ifs, ",")
	case "hccl", "oneccl":
		if flavor.Vendor != "intel" {
			return CollectiveEnvPlan{}, fmt.Errorf("%s requires intel flavor", flavor.Capabilities.Collective)
		}
		env["CCL_WORKER_COUNT"] = "1"
		if len(ifs) > 0 {
			env["CCL_ATL_TRANSPORT"] = "ofi"
		}
	default:
		return CollectiveEnvPlan{}, fmt.Errorf("unsupported collective %q", flavor.Capabilities.Collective)
	}
	return CollectiveEnvPlan{Backend: flavor.Capabilities.Collective, Env: env}, nil
}

func fabricNames(records []clusterprobe.HostRecord) []string {
	seen := map[string]bool{}
	var names []string
	for _, r := range records {
		for _, f := range r.Fabric {
			name := f.Name
			if f.RDMADevice != "" {
				name = f.RDMADevice
			}
			if name != "" && !seen[name] {
				seen[name] = true
				names = append(names, name)
			}
		}
	}
	return names
}
