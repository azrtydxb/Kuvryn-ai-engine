package clusterprobe

import "time"

type RDMAValidation struct {
	Required     bool      `json:"required"`
	Reachable    bool      `json:"reachable"`
	TransferOK   bool      `json:"transferOk"`
	CollectiveOK bool      `json:"collectiveOk"`
	PodVisible   bool      `json:"podVisible"`
	CheckedAt    time.Time `json:"checkedAt"`
	Message      string    `json:"message,omitempty"`
}

func (v RDMAValidation) UsableForRequired(now time.Time, maxAge time.Duration, kubernetes bool) bool {
	if !v.Required {
		return true
	}
	if v.CheckedAt.IsZero() || now.Sub(v.CheckedAt) > maxAge {
		return false
	}
	if kubernetes && !v.PodVisible {
		return false
	}
	return v.Reachable && v.TransferOK && v.CollectiveOK
}
