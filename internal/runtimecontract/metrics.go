package runtimecontract

import (
	"fmt"
	"io"
	"time"
)

type Metrics struct {
	Flavor       string
	StartedAt    time.Time
	Ready        bool
	RestartCount int
	OOMCount     int
}

func WritePrometheus(w io.Writer, m Metrics) error {
	uptime := 0.0
	if !m.StartedAt.IsZero() {
		uptime = time.Since(m.StartedAt).Seconds()
	}
	ready := 0
	if m.Ready {
		ready = 1
	}
	_, err := fmt.Fprintf(w, "kuvryn_runtime_uptime_seconds{flavor=%q} %.0f\nkuvryn_runtime_ready{flavor=%q} %d\nkuvryn_runtime_restarts_total{flavor=%q} %d\nkuvryn_runtime_oom_total{flavor=%q} %d\n", m.Flavor, uptime, m.Flavor, ready, m.Flavor, m.RestartCount, m.Flavor, m.OOMCount)
	return err
}
