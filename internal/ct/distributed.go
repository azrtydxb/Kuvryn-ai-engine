package ct

import "fmt"

type DistributedResult struct {
	Flavor     string `json:"flavor"`
	WorldSize  int    `json:"worldSize"`
	RanksReady int    `json:"ranksReady"`
	Backend    string `json:"backend"`
	Inference  bool   `json:"inference"`
}

func ValidateDistributedResult(r DistributedResult) error {
	if r.WorldSize < 2 {
		return fmt.Errorf("distributed CT requires at least two ranks")
	}
	if r.RanksReady != r.WorldSize {
		return fmt.Errorf("ranks ready %d does not match world size %d", r.RanksReady, r.WorldSize)
	}
	if r.Backend == "" {
		return fmt.Errorf("collective backend is required")
	}
	if !r.Inference {
		return fmt.Errorf("distributed CT must include one inference request")
	}
	return nil
}
