package aggregate

import "github.com/ecce-machina/machina-trace/internal/diff"

type NodeObservation struct {
	Node    string
	StartNS int64
	EndNS   int64
	Deltas  []diff.CounterDelta
}

func NewNodeObservation(deltas []diff.CounterDelta) (NodeObservation, bool) {
	if len(deltas) == 0 {
		return NodeObservation{}, false
	}

	observation := NodeObservation{
		Node:    deltas[0].Node,
		StartNS: deltas[0].StartNS,
		EndNS:   deltas[0].EndNS,
		Deltas:  deltas,
	}

	for _, delta := range deltas[1:] {
		if delta.Node != observation.Node {
			return NodeObservation{}, false
		}

		if delta.StartNS < observation.StartNS {
			observation.StartNS = delta.StartNS
		}

		if delta.EndNS > observation.EndNS {
			observation.EndNS = delta.EndNS
		}
	}

	return observation, true
}
