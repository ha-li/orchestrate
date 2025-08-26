package task

// These are the valid state transitions, represented as a map.
// The key is the current state, the value is an array of the possible future states
// when starting from the current state.
// if the state transition presented by (key=current state, value=next state)
// does not exist, then the state transition is invalid
var stateTransitionMap = map[State][]State{
	Pending:   []State{Scheduled},
	Scheduled: []State{Scheduled, Running, Failed},
	Running:   []State{Running, Completed, Failed},
	Completed: []State{}, // completed is the end state
	Failed:    []State{}, // failed is an end state
}

// returns true if the src, dst is a valid state transition
// otherwise returns false
func ValidStateTransition(src State, dst State) bool {
	return contains(stateTransitionMap[src], dst)
}

func contains(states []State, state State) bool {
	for _, s := range states {
		if s == state {
			return true
		}
	}
	return false
}
