package placement

import "github.com/azrtydxb/kuvryn-ai-engine/internal/engineimage"

func CompatibleFlavor(req DeploymentRequest, flavor engineimage.ImageManifest) bool {
	if req.Engine != flavor.Engine || req.Vendor != flavor.Vendor {
		return false
	}
	if req.Distributed && !flavor.Capabilities.Distributed {
		return false
	}
	return true
}
