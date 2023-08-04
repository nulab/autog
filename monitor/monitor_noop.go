package monitor

func NewNoop() Monitor {
	return &noopMonitor{}
}

type noopMonitor struct{}

func (m *noopMonitor) Send(string, any) {
	return
}

func (m *noopMonitor) Close() {
	return
}
