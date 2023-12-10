package goshmuffle

type state uint8

const (
	stateRunning state = iota + 1
	stateDone
)

func (r *runner) setStateRunning() {
	r.mu.Lock()
	r.state = stateRunning
	r.mu.Unlock()
}

func (r *runner) setStateDone() {
	r.mu.Lock()
	r.state = stateDone
	r.mu.Unlock()
}

func (r *runner) getState() state {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.state
}

func (r *runner) isDone() bool {
	return r.getState() == stateDone
}

func (r *runner) isRunning() bool {
	return r.getState() == stateRunning
}
