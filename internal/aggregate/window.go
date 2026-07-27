package aggregate

import (
	"errors"

	"github.com/ecce-machina/machina-trace/internal/diff"
    "github.com/ecce-machina/machina-trace/internal/features"
)

var (
	ErrObservationOutsideWindow = errors.New("observation falls outside cluster window")
	ErrDuplicateNode            = errors.New("cluster window already contains observation for node")
	ErrEmptyNode                = errors.New("observation node is empty")
)

type ClusterWindow struct {
	StartNS int64
	EndNS   int64
	Nodes   map[string]NodeObservation
}

func NewClusterWindow(startNS, endNS int64) ClusterWindow {
	return ClusterWindow{
		StartNS: startNS,
		EndNS:   endNS,
		Nodes:   make(map[string]NodeObservation),
	}
}

func (w *ClusterWindow) Add(obs NodeObservation) error {
	if obs.Node == "" {
		return ErrEmptyNode
	}

	if obs.StartNS < w.StartNS || obs.EndNS > w.EndNS {
		return ErrObservationOutsideWindow
	}

	if _, exists := w.Nodes[obs.Node]; exists {
		return ErrDuplicateNode
	}

	w.Nodes[obs.Node] = obs
	return nil
}

func (w ClusterWindow) Node(name string) (NodeObservation, bool) {
	obs, ok := w.Nodes[name]
	return obs, ok
}

func (w ClusterWindow) NumNodes() int {
	return len(w.Nodes)
}

func (w ClusterWindow) Deltas() []diff.CounterDelta {
	var total int

	for _, obs := range w.Nodes {
		total += len(obs.Deltas)
	}

	out := make([]diff.CounterDelta, 0, total)

	for _, obs := range w.Nodes {
		out = append(out, obs.Deltas...)
	}

	return out
}

// WorkloadFeatures derives per-node workload features and the corresponding
// filesystem-wide aggregate for this cluster window.
func (w ClusterWindow) WorkloadFeatures() FilesystemFeatures {
	nodes := make(map[string]features.WorkloadFeatures, len(w.Nodes))

	for node, observation := range w.Nodes {
		nodes[node] = features.WorkloadFeaturesFromDeltas(
			observation.Deltas,
		)
	}

	return AggregateWorkloadFeatures(nodes)
}
