package monitor

type FilterFn func(phase int, alg, key string) bool

func MatchAll(phase int, alg, key string) FilterFn {
	return func(p int, a, k string) bool {
		return p == phase && a == alg && k == key
	}
}

type filteredMonitor struct {
	m Monitor
	f FilterFn
}

func (m filteredMonitor) Log(phase int, alg, key string, val any) {
	if m.f(phase, alg, key) {
		m.m.Log(phase, alg, key, val)
	}
}

func WrapFilter(m Monitor, f FilterFn) Monitor {
	return filteredMonitor{m, f}
}
