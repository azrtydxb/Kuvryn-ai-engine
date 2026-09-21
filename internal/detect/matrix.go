package detect

import "github.com/azrtydxb/kuvryn-ai-engine/internal/engineimage"

const UpdatePlanSchema = "kuvryn.engine-update-plan/v1"

type UpdatePlan struct {
	Schema  string       `json:"schema"`
	Flavors []FlavorPlan `json:"flavors"`
	Errors  []PlanError  `json:"errors,omitempty"`
}

type FlavorPlan struct {
	Name           string                        `json:"name"`
	Engine         string                        `json:"engine"`
	Vendor         string                        `json:"vendor"`
	Manifest       string                        `json:"manifest"`
	Dockerfile     string                        `json:"dockerfile"`
	Image          string                        `json:"image"`
	Certification  string                        `json:"certification"`
	Changed        bool                          `json:"changed"`
	Reason         string                        `json:"reason"`
	ResolvedInputs []engineimage.ResolvedInput   `json:"resolvedInputs"`
	Tests          []engineimage.TestDeclaration `json:"tests"`
}

type PlanError struct {
	Manifest string `json:"manifest,omitempty"`
	Flavor   string `json:"flavor,omitempty"`
	Error    string `json:"error"`
}
