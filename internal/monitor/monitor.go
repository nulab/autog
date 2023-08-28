package monitor

import "github.com/nulab/autog/internal/processor"

type Monitor interface {
	Log(phase int, alg, key string, val any)
}

// globals
var (
	m Monitor
	p int
	a string
)

func Set(monitor Monitor) {
	if monitor != nil {
		m = monitor
	}
}

func PrefixFor(proc processor.P) {
	if m != nil {
		p = proc.Phase()
		a = proc.String()
	}
}

func Reset() {
	if m != nil {
		m = nil
		p = 0
		a = ""
	}
}

func Log(key string, val any) {
	if m != nil {
		m.Log(p, a, key, val)
	}
}
