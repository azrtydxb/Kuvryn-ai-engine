package ops

type FailureKind string

const (
	FailureRuntimeCrash  FailureKind = "runtime-crash"
	FailureDriverCrash   FailureKind = "driver-crash"
	FailureCgroupOOM     FailureKind = "cgroup-oom"
	FailureHostOOM       FailureKind = "host-oom"
	FailureEvicted       FailureKind = "kubernetes-eviction"
	FailureHealthRestart FailureKind = "health-check-restart"
	FailureUserStop      FailureKind = "user-stop"
	FailureUnknown       FailureKind = "unknown"
)

type FailureEvidence struct {
	Kind    FailureKind `json:"kind"`
	Source  string      `json:"source"`
	Message string      `json:"message,omitempty"`
}

func ClassifyContainerExit(exitCode int, oomKilled bool, userStopped bool) FailureEvidence {
	if userStopped {
		return FailureEvidence{Kind: FailureUserStop, Source: "container"}
	}
	if oomKilled || exitCode == 137 {
		return FailureEvidence{Kind: FailureCgroupOOM, Source: "container"}
	}
	if exitCode != 0 {
		return FailureEvidence{Kind: FailureRuntimeCrash, Source: "container"}
	}
	return FailureEvidence{Kind: FailureUnknown, Source: "container"}
}

func ClassifyKubernetes(reason string) FailureEvidence {
	switch reason {
	case "OOMKilled":
		return FailureEvidence{Kind: FailureCgroupOOM, Source: "kubernetes"}
	case "SystemOOM":
		return FailureEvidence{Kind: FailureHostOOM, Source: "kubernetes"}
	case "Evicted":
		return FailureEvidence{Kind: FailureEvicted, Source: "kubernetes"}
	case "Unhealthy":
		return FailureEvidence{Kind: FailureHealthRestart, Source: "kubernetes"}
	default:
		return FailureEvidence{Kind: FailureUnknown, Source: "kubernetes", Message: reason}
	}
}
