package monitor

type funcMonitor func(phase int, alg, key string, val any)

func (f funcMonitor) Log(phase int, alg, key string, val any) {
	f(phase, alg, key, val)
}

func NewFunc(fn func(phase int, alg, key string, val any)) Monitor {
	return funcMonitor(fn)
}
