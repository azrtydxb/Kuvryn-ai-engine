package clusterprobe

type TopologyClass string

const (
	TopologySingleNode       TopologyClass = "single-node"
	TopologyDirect           TopologyClass = "direct"
	TopologyRing             TopologyClass = "ring"
	TopologySwitch           TopologyClass = "switch"
	TopologyKubernetesFabric TopologyClass = "kubernetes-fabric"
	TopologyEthernetOnly     TopologyClass = "ethernet-only"
	TopologyUnknown          TopologyClass = "unknown"
)

func ClassifyTopology(records []HostRecord, kubernetes bool) TopologyClass {
	if len(records) == 0 {
		return TopologyUnknown
	}
	if len(records) == 1 {
		return TopologySingleNode
	}
	withFabric := 0
	podVisible := 0
	for _, r := range records {
		if len(r.Fabric) > 0 {
			withFabric++
		}
		for _, f := range r.Fabric {
			if f.PodVisible {
				podVisible++
			}
		}
	}
	if withFabric == 0 {
		return TopologyEthernetOnly
	}
	if kubernetes && podVisible > 0 {
		return TopologyKubernetesFabric
	}
	if withFabric == len(records) && len(records) == 2 {
		return TopologyDirect
	}
	if withFabric == len(records) {
		return TopologySwitch
	}
	return TopologyUnknown
}
