package aggregate

import (
	"errors"
	"testing"

	"github.com/ecce-machina/machina-trace/internal/diff"
)

func TestClusterWindowAdd(t *testing.T) {
	window := NewClusterWindow(
		1_000_000_000,
		6_000_000_000,
	)

	obs := NodeObservation{
		Node:    "node01",
		StartNS: 2_000_000_000,
		EndNS:   5_000_000_000,
		Deltas: []diff.CounterDelta{
			{
				Node:      "node01",
				Collector: "proc_diskstats",
			},
			{
				Node:      "node01",
				Collector: "proc_net_dev",
			},
		},
	}

	if err := window.Add(obs); err != nil {
		t.Fatalf("unexpected error adding observation: %v", err)
	}

	if window.NumNodes() != 1 {
		t.Fatalf("expected 1 node, got %d", window.NumNodes())
	}

	got, ok := window.Node("node01")
	if !ok {
		t.Fatal("expected to find node01")
	}

	if got.Node != "node01" {
		t.Errorf("expected node01, got %q", got.Node)
	}

	if len(window.Deltas()) != 2 {
		t.Errorf("expected 2 flattened deltas, got %d", len(window.Deltas()))
	}
}

func TestClusterWindowRejectsDuplicateNode(t *testing.T) {
	window := NewClusterWindow(1, 10)

	obs := NodeObservation{
		Node:    "node01",
		StartNS: 2,
		EndNS:   9,
	}

	if err := window.Add(obs); err != nil {
		t.Fatalf("unexpected first add error: %v", err)
	}

	err := window.Add(obs)

	if !errors.Is(err, ErrDuplicateNode) {
		t.Fatalf("expected ErrDuplicateNode, got %v", err)
	}
}

func TestClusterWindowRejectsObservationBeforeWindow(t *testing.T) {
	window := NewClusterWindow(10, 20)

	obs := NodeObservation{
		Node:    "node01",
		StartNS: 9,
		EndNS:   15,
	}

	err := window.Add(obs)

	if !errors.Is(err, ErrObservationOutsideWindow) {
		t.Fatalf(
			"expected ErrObservationOutsideWindow, got %v",
			err,
		)
	}
}

func TestClusterWindowRejectsObservationAfterWindow(t *testing.T) {
	window := NewClusterWindow(10, 20)

	obs := NodeObservation{
		Node:    "node01",
		StartNS: 15,
		EndNS:   21,
	}

	err := window.Add(obs)

	if !errors.Is(err, ErrObservationOutsideWindow) {
		t.Fatalf(
			"expected ErrObservationOutsideWindow, got %v",
			err,
		)
	}
}

func TestClusterWindowRejectsEmptyNode(t *testing.T) {
	window := NewClusterWindow(1, 10)

	obs := NodeObservation{
		StartNS: 2,
		EndNS:   9,
	}

	err := window.Add(obs)

	if !errors.Is(err, ErrEmptyNode) {
		t.Fatalf("expected ErrEmptyNode, got %v", err)
	}
}

func TestClusterWindowWorkloadFeaturesPreservesNodes(t *testing.T) {
	window := NewClusterWindow(1, 10)

	err := window.Add(NodeObservation{
		Node:    "client-01",
		StartNS: 1,
		EndNS:   10,
	})
	if err != nil {
		t.Fatalf("add client-01: %v", err)
	}

	err = window.Add(NodeObservation{
		Node:    "client-02",
		StartNS: 1,
		EndNS:   10,
	})
	if err != nil {
		t.Fatalf("add client-02: %v", err)
	}

	got := window.WorkloadFeatures()

	if len(got.Nodes) != 2 {
		t.Fatalf("got %d nodes, want 2", len(got.Nodes))
	}

	if _, ok := got.Nodes["client-01"]; !ok {
		t.Fatal("client-01 features missing")
	}

	if _, ok := got.Nodes["client-02"]; !ok {
		t.Fatal("client-02 features missing")
	}
}
