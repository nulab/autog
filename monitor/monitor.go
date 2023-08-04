package monitor

type Monitor interface {
	Send(key string, val any)
	Close()
}

type Log struct {
	Name  string
	Value value
}
