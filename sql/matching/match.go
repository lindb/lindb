package matching

type Match struct {
	Captures *Captures
}

func NewMatch(captures *Captures) *Match {
	return &Match{
		Captures: captures,
	}
}

func (m *Match) Capture(capture *Capture) any {
	return m.Captures.Get(capture)
}
