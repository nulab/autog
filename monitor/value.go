package monitor

type value struct{ any }

func (v value) AsInt() int {
	return v.any.(int)
}
