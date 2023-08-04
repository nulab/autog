package monitor

func New(c chan Log) Monitor {
	return &monitor{c}
}

type monitor struct {
	c chan Log
}

func (m *monitor) Send(key string, val any) {
	if m == nil {
		return
	}
	m.c <- Log{key, value{val}}
}

func (m *monitor) Close() {
	if m == nil {
		return
	}
	close(m.c)
}
