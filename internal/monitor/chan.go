package monitor

type chanMonitor struct {
	c chan any
}

func (m *chanMonitor) Log(_ int, _, _ string, val any) {
	m.c <- val
}

func NewFilteredChan(c chan any, filter FilterFn) Monitor {
	return WrapFilter(&chanMonitor{c}, filter)
}
