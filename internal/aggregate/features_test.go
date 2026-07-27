package aggregate

import (
	"testing"

	"github.com/ecce-machina/machina-trace/internal/features"
)

func TestAggregateWorkloadFeatures(t *testing.T) {
	nodes := map[string]features.WorkloadFeatures{
		"client-01": {
			IntervalSec:       10,
			ReadOpsPerSec:     10,
			WriteOpsPerSec:    5,
			ReadBytesPerSec:   40 * 1024,
			WriteBytesPerSec:  40 * 1024,
			MetadataOpsPerSec: 12,
			LockOpsPerSec:     4,
			MDCRPCsPerSec:     8,
		},
		"client-02": {
			IntervalSec:       10,
			ReadOpsPerSec:     30,
			WriteOpsPerSec:    15,
			ReadBytesPerSec:   240 * 1024,
			WriteBytesPerSec:  240 * 1024,
			MetadataOpsPerSec: 20,
			LockOpsPerSec:     6,
			MDCRPCsPerSec:     14,
		},
	}

	got := AggregateWorkloadFeatures(nodes)

	if len(got.Nodes) != 2 {
		t.Fatalf("got %d nodes, want 2", len(got.Nodes))
	}

	assertFloat(t, "IntervalSec", got.Aggregate.IntervalSec, 10)

	assertFloat(t, "ReadOpsPerSec", got.Aggregate.ReadOpsPerSec, 40)
	assertFloat(t, "WriteOpsPerSec", got.Aggregate.WriteOpsPerSec, 20)

	assertFloat(
		t,
		"ReadBytesPerSec",
		got.Aggregate.ReadBytesPerSec,
		280*1024,
	)
	assertFloat(
		t,
		"WriteBytesPerSec",
		got.Aggregate.WriteBytesPerSec,
		280*1024,
	)

	// 280 KiB/s divided by 40 operations/s.
	assertFloat(
		t,
		"AverageReadSizeBytes",
		got.Aggregate.AverageReadSizeBytes,
		7*1024,
	)

	// 280 KiB/s divided by 20 operations/s.
	assertFloat(
		t,
		"AverageWriteSizeBytes",
		got.Aggregate.AverageWriteSizeBytes,
		14*1024,
	)

	assertFloat(
		t,
		"ReadFraction",
		got.Aggregate.ReadFraction,
		2.0/3.0,
	)
	assertFloat(
		t,
		"WriteFraction",
		got.Aggregate.WriteFraction,
		1.0/3.0,
	)

	assertFloat(
		t,
		"MetadataOpsPerSec",
		got.Aggregate.MetadataOpsPerSec,
		32,
	)
	assertFloat(
		t,
		"LockOpsPerSec",
		got.Aggregate.LockOpsPerSec,
		10,
	)
	assertFloat(
		t,
		"MDCRPCsPerSec",
		got.Aggregate.MDCRPCsPerSec,
		22,
	)
}

func TestAggregateWorkloadFeaturesDoesNotAverageDerivedValues(t *testing.T) {
	nodes := map[string]features.WorkloadFeatures{
		"small-reader": {
			IntervalSec:          10,
			ReadOpsPerSec:        100,
			ReadBytesPerSec:      100 * 1024,
			AverageReadSizeBytes: 1024,
		},
		"large-reader": {
			IntervalSec:          10,
			ReadOpsPerSec:        1,
			ReadBytesPerSec:      1024 * 1024,
			AverageReadSizeBytes: 1024 * 1024,
		},
	}

	got := AggregateWorkloadFeatures(nodes)

	want := float64((100 * 1024) + (1024 * 1024))
	want /= 101

	assertFloat(
		t,
		"AverageReadSizeBytes",
		got.Aggregate.AverageReadSizeBytes,
		want,
	)
}

func TestAggregateWorkloadFeaturesPreservesNodeFeatures(t *testing.T) {
	input := map[string]features.WorkloadFeatures{
		"client-01": {
			IntervalSec:     5,
			ReadOpsPerSec:   7,
			WriteOpsPerSec:  3,
			ReadBytesPerSec: 4096,
		},
	}

	got := AggregateWorkloadFeatures(input)

	node, ok := got.Nodes["client-01"]
	if !ok {
		t.Fatal("client-01 was not preserved")
	}

	assertFloat(t, "node ReadOpsPerSec", node.ReadOpsPerSec, 7)
	assertFloat(t, "node WriteOpsPerSec", node.WriteOpsPerSec, 3)
	assertFloat(t, "node ReadBytesPerSec", node.ReadBytesPerSec, 4096)

	// The result owns its map. Changing the input map must not remove the
	// preserved node from the result.
	delete(input, "client-01")

	if _, ok := got.Nodes["client-01"]; !ok {
		t.Fatal("result shares its node map with the caller")
	}
}

func TestAggregateWorkloadFeaturesEmptyInput(t *testing.T) {
	got := AggregateWorkloadFeatures(nil)

	if len(got.Nodes) != 0 {
		t.Fatalf("got %d nodes, want 0", len(got.Nodes))
	}

	want := features.WorkloadFeatures{}
	if got.Aggregate != want {
		t.Fatalf("got aggregate %+v, want zero value", got.Aggregate)
	}
}

func TestAggregateWorkloadFeaturesZeroOperations(t *testing.T) {
	nodes := map[string]features.WorkloadFeatures{
		"client-01": {
			IntervalSec:       10,
			MetadataOpsPerSec: 5,
		},
	}

	got := AggregateWorkloadFeatures(nodes)

	assertFloat(
		t,
		"AverageReadSizeBytes",
		got.Aggregate.AverageReadSizeBytes,
		0,
	)
	assertFloat(
		t,
		"AverageWriteSizeBytes",
		got.Aggregate.AverageWriteSizeBytes,
		0,
	)
	assertFloat(t, "ReadFraction", got.Aggregate.ReadFraction, 0)
	assertFloat(t, "WriteFraction", got.Aggregate.WriteFraction, 0)
}

func assertFloat(t *testing.T, name string, got, want float64) {
	t.Helper()

	const tolerance = 1e-9

	difference := got - want
	if difference < 0 {
		difference = -difference
	}

	if difference > tolerance {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
