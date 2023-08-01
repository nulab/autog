package monitor

type Log struct {
	Name  string
	Value value
}

func New(c chan Log) *Monitor {
	return &Monitor{c}
}

type Monitor struct {
	c chan Log
}

func (m *Monitor) Send(key string, val any) {
	if m == nil {
		return
	}
	m.c <- Log{key, value{val}}
}

func (m *Monitor) Close() {
	if m == nil {
		return
	}
	close(m.c)
}
