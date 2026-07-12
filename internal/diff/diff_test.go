package diff

import (
	"testing"

	"github.com/ecce-machina/machina-trace/internal/snapshot"
)

func TestDiffSnapshotsPreservesSourceContext(t *testing.T) {
	older := &snapshot.Snapshot{
		Node:        "node01",
		TimestampNS: 1_000_000_000,
		Sources: []snapshot.Source{
			{
				Collector:   "proc_diskstats",
				Object:      "sda",
				Mount:       "/lustre",
				TimestampNS: 1_000_000_000,
				Counters: map[string]int64{
					"reads_completed": 100,
					"sectors_read":    200,
				},
			},
		},
	}

	newer := &snapshot.Snapshot{
		Node:        "node01",
		TimestampNS: 3_000_000_000,
		Sources: []snapshot.Source{
			{
				Collector:   "proc_diskstats",
				Object:      "sda",
				Mount:       "/lustre",
				TimestampNS: 3_000_000_000,
				Counters: map[string]int64{
					"reads_completed": 110,
					"sectors_read":    240,
				},
			},
		},
	}

	deltas := DiffSnapshots(older, newer)

	if len(deltas) != 1 {
		t.Fatalf("expected 1 delta, got %d", len(deltas))
	}

	got := deltas[0]

	if got.Node != "node01" {
		t.Errorf("expected node01, got %q", got.Node)
	}

	if got.Collector != "proc_diskstats" {
		t.Errorf("expected proc_diskstats, got %q", got.Collector)
	}

	if got.Object != "sda" {
		t.Errorf("expected sda, got %q", got.Object)
	}

	if got.Mount != "/lustre" {
		t.Errorf("expected /lustre, got %q", got.Mount)
	}

	if got.StartNS != 1_000_000_000 {
		t.Errorf("expected start 1000000000, got %d", got.StartNS)
	}

	if got.EndNS != 3_000_000_000 {
		t.Errorf("expected end 3000000000, got %d", got.EndNS)
	}

	if got.IntervalSec != 2 {
		t.Errorf("expected interval 2, got %f", got.IntervalSec)
	}

	if got.Deltas["reads_completed"] != 10 {
		t.Errorf(
			"expected reads_completed delta 10, got %d",
			got.Deltas["reads_completed"],
		)
	}

	if got.Rates["reads_completed"] != 5 {
		t.Errorf(
			"expected reads_completed rate 5, got %f",
			got.Rates["reads_completed"],
		)
	}
}
