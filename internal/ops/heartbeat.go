package ops

import "time"

type Heartbeat struct {
	DeploymentID string    `json:"deploymentId"`
	ImageDigest  string    `json:"imageDigest"`
	Flavor       string    `json:"flavor"`
	Node         string    `json:"node,omitempty"`
	Pod          string    `json:"pod,omitempty"`
	Container    string    `json:"container,omitempty"`
	State        string    `json:"state"`
	Ready        bool      `json:"ready"`
	Rank         int       `json:"rank,omitempty"`
	WorldSize    int       `json:"worldSize,omitempty"`
	At           time.Time `json:"at"`
}

func NewHeartbeat(deploymentID, digest, flavor, state string, ready bool) Heartbeat {
	return Heartbeat{DeploymentID: deploymentID, ImageDigest: digest, Flavor: flavor, State: state, Ready: ready, At: time.Now().UTC()}
}
