package ordering

import (
	"github.com/nulab/autog/monitor"
)

type phase3monitor struct {
	alg string
	m   *monitor.Monitor
}

func (m phase3monitor) Monitor(key string, value any) {
	if m.m != nil {
		m.m.Send("phase3/"+m.alg+"/"+key, value)
	}
}
