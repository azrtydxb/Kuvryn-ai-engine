package runtimecontract

import (
	"encoding/json"
	"os"
	"time"
)

type State string

const (
	StateStarting State = "starting"
	StateReady    State = "ready"
	StateStopping State = "stopping"
	StateExited   State = "exited"
	StateFailed   State = "failed"
)

type RuntimeState struct {
	State     State     `json:"state"`
	PID       int       `json:"pid,omitempty"`
	Ready     bool      `json:"ready"`
	UpdatedAt time.Time `json:"updatedAt"`
	Message   string    `json:"message,omitempty"`
}

type RuntimeEvent struct {
	State   State     `json:"state"`
	PID     int       `json:"pid,omitempty"`
	Ready   bool      `json:"ready"`
	At      time.Time `json:"at"`
	Message string    `json:"message,omitempty"`
}

func WriteState(path string, state RuntimeState) error {
	state.UpdatedAt = time.Now().UTC()
	b, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

func AppendEvent(path string, event RuntimeEvent) error {
	event.At = time.Now().UTC()
	b, err := json.Marshal(event)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(append(b, '\n'))
	return err
}
