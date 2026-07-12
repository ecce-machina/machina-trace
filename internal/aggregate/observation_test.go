package aggregate

import (
	"testing"

	"github.com/ecce-machina/machina-trace/internal/diff"
)

func TestNewNodeObservation(t *testing.T) {
	deltas := []diff.CounterDelta{
		{
			Node:      "node01",
			Collector: "proc_diskstats",
			StartNS:   2_000_000_000,
			EndNS:     4_000_000_000,
		},
		{
			Node:      "node01",
			Collector: "proc_net_dev",
			StartNS:   1_000_000_000,
			EndNS:     5_000_000_000,
		},
	}

	got, ok := NewNodeObservation(deltas)
	if !ok {
		t.Fatal("expected observation to be created")
	}

	if got.Node != "node01" {
		t.Errorf("expected node01, got %q", got.Node)
	}

	if got.StartNS != 1_000_000_000 {
		t.Errorf("expected earliest start, got %d", got.StartNS)
	}

	if got.EndNS != 5_000_000_000 {
		t.Errorf("expected latest end, got %d", got.EndNS)
	}

	if len(got.Deltas) != 2 {
		t.Errorf("expected 2 deltas, got %d", len(got.Deltas))
	}
}

func TestNewNodeObservationRejectsMultipleNodes(t *testing.T) {
	deltas := []diff.CounterDelta{
		{
			Node:    "node01",
			StartNS: 1,
			EndNS:   2,
		},
		{
			Node:    "node02",
			StartNS: 1,
			EndNS:   2,
		},
	}

	_, ok := NewNodeObservation(deltas)
	if ok {
		t.Fatal("expected mixed-node observation to be rejected")
	}
}

func TestNewNodeObservationRejectsEmptyInput(t *testing.T) {
	_, ok := NewNodeObservation(nil)
	if ok {
		t.Fatal("expected empty observation to be rejected")
	}
}
