package runtimecontract

import "net/http"

type Health struct {
	Live    func() bool
	Ready   func() bool
	Startup func() bool
}

func Handler(h Health) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/kuvryn/health/live", status(h.Live))
	mux.HandleFunc("/kuvryn/health/ready", status(h.Ready))
	mux.HandleFunc("/kuvryn/health/startup", status(h.Startup))
	return mux
}

func status(fn func() bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if fn == nil || !fn() {
			http.Error(w, "not ready", http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte("ok\n"))
	}
}
