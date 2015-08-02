package antagonist

// Params ...
type Params []string

func (p *Params) index(param string) int {
	for i, pp := range *p {
		if pp == param {
			return i
		}
	}

	return -1
}

// Set ...
func (p *Params) Set(param string) {
	i := p.index(param)
	if i != -1 {
		(*p)[i] = param

		return
	}

	*p = append(*p, param)
}

// Del ...
func (p *Params) Del(param string) {
	i := p.index(param)
	if i == -1 {
		return
	}

	*p = append((*p)[:i], (*p)[i+1:]...)
}

// Exists ...
func (p *Params) Exists(param string) bool {
	return p.index(param) == -1
}
