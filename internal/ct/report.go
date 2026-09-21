package ct

import (
	"encoding/json"
	"errors"
	"time"
)

const ReportSchema = "kuvryn.container-test-report/v1"

type Report struct {
	Schema    string    `json:"schema"`
	Flavor    string    `json:"flavor"`
	Image     string    `json:"image"`
	Digest    string    `json:"digest,omitempty"`
	StartedAt time.Time `json:"startedAt"`
	EndedAt   time.Time `json:"endedAt"`
	Results   []Result  `json:"results"`
}

type Result struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Status   string `json:"status"`
	Required bool   `json:"required"`
	Message  string `json:"message,omitempty"`
}

func (r Report) Validate() error {
	if r.Schema != ReportSchema {
		return errors.New("invalid report schema")
	}
	if r.Flavor == "" || r.Image == "" {
		return errors.New("flavor and image are required")
	}
	if len(r.Results) == 0 {
		return errors.New("at least one result is required")
	}
	for _, res := range r.Results {
		if res.Name == "" || res.Type == "" || res.Status == "" {
			return errors.New("result name, type, and status are required")
		}
	}
	return nil
}

func (r Report) JSON() ([]byte, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	return json.MarshalIndent(r, "", "  ")
}
