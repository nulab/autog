package phase3

import (
	"github.com/nulab/autog/monitor"
)

type phase3monitor struct {
	alg string
	monitor.Monitor
}

func (m phase3monitor) Send(key string, value any) {
	if m.Monitor != nil {
		m.Monitor.Send("phase3/"+m.alg+"/"+key, value)
	}
}
